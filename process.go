package proc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

const startTimeFieldNum = 21

type ProcStatusFunc func(pid int) (map[string]string, error)

// Process stores static information, like process ID and name.
// Dynamic information, like memory usage or uptime, should be read from /proc
// every time they're needed to avoid worrying about stale information.
type Process struct {
	Pid       int      // process ID
	Command   string   // command name
	Args      []string // command line arguments
	StartTime uint64   // process start time in clock ticks since system boot

	status ProcStatusFunc // function to call for proc status info (used internally in proc pkg)
}

// status returns the information from /proc/[pid]/status in a map.
func status(pid int) (map[string]string, error) {
	// get the proc's status
	statusBytes, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "status"))
	if err != nil {
		log.Print("Error reading proc status (process might have exited)", err)
		return nil, err
	}

	// split each line on ":" and build a map of attribute name -> value
	status := make(map[string]string)
	trimmedStatus := strings.TrimSpace(string(statusBytes))
	for _, line := range strings.Split(trimmedStatus, "\n") {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		status[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
	}

	return status, nil
}

// parseMemorySize takes a string from /proc/[pid]/status, verifies that it is
// of the form "1353584 kB", and returns the number of *bytes* expressed by that
// string as a uint64 (in this case, 1386070016).
func parseMemorySize(memStr string) (uint64, error) {
	fields := strings.Fields(memStr)

	if len(fields) != 2 {
		return 0, errors.New("memory size string not of form '1353584 kB'")
	}

	sizeStr := strings.TrimSpace(fields[0])
	sizeNum, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	units := strings.TrimSpace(fields[1])
	var sizeBytes uint64 = 0
	if units == "kB" {
		sizeBytes = sizeNum * 1024
	} else if units == "B" {
		sizeBytes = sizeNum
	} else {
		return 0, errors.New(fmt.Sprintf("unrecognized units: %v", units))
	}

	return sizeBytes, nil
}

// MemUsage returns a process' current memory usage:
// vmSize = total allocated memory (RAM, swap, disk, etc.) in bytes. TODO: This seems to include mem in huge pages?
// rssSize = memory currently resident in RAM in bytes
// hugeSize = memory in huge pages in bytes
func (p *Process) MemUsage() (vmSizeBytes, rssSizeBytes, hugeSizeBytes uint64, err error) {
	// return 0 if /proc/[pid]/status doesn't have info for that key
	vmSizeBytes, rssSizeBytes, hugeSizeBytes = 0, 0, 0

	// proc status info as a map
	status, err := p.status(p.Pid)
	if err != nil {
		return
	}

	vmSizeStr, ok := status["VmSize"]
	if ok {
		vmSizeBytes, err = parseMemorySize(vmSizeStr)
		if err != nil {
			log.Printf("Error parsing VmSize for %v: %v", p.Command, vmSizeStr)
			vmSizeBytes = 0
		}
	}

	rssSizeStr, ok := status["VmRSS"]
	if ok {
		rssSizeBytes, err = parseMemorySize(rssSizeStr)
		if err != nil {
			log.Printf("Error parsing VmRSS for %v: %v", p.Command, rssSizeStr)
			rssSizeBytes = 0
		}
	}

	hugeSizeStr, ok := status["HugetlbPages"]
	if ok {
		hugeSizeBytes, err = parseMemorySize(hugeSizeStr)
		if err != nil {
			log.Printf("Error parsing HugetlbPages for %v: %v", p.Command, hugeSizeStr)
			hugeSizeBytes = 0
		}
	}

	return vmSizeBytes, rssSizeBytes, hugeSizeBytes, nil
}

// State returns the state of the process, which is one of:
//	R  Running
//	S  Sleeping in an interruptible wait
//	D  Waiting in uninterruptible disk sleep
//	Z  Zombie
//	T  Stopped (on a signal) or (before Linux 2.6.33) trace stopped
//	t  Tracing stop (Linux 2.6.33 onward)
//	W  Paging (only before Linux 2.6.0)
//	X  Dead (from Linux 2.6.0 onward)
//	x  Dead (Linux 2.6.33 to 3.13 only)
//	K  Wakekill (Linux 2.6.33 to 3.13 only)
//	W  Waking (Linux 2.6.33 to 3.13 only)
//	P  Parked (Linux 3.9 to 3.13 only)
//
// status: one of the character codes above
// description: a short description of the state, e.g., "sleeping"
func (p *Process) State() (state string, description string, err error) {
	// proc status info as a map
	status, err := p.status(p.Pid)
	if err != nil {
		return "", "", err
	}

	stateStr, ok := status["State"]
	if !ok {
		return "", "", errors.New("'State' not found in process info")
	}

	// stateStr should be of the form "S (sleeping)"
	fields := strings.Fields(stateStr)
	if len(fields) != 2 {
		return "", "", errors.New(fmt.Sprintf("'State' not in expected format: %v", stateStr))
	}

	return fields[0], strings.Trim(fields[1], "()"), nil
}

// Processes returns a list of Process structs, one for each running process
func Processes() ([]Process, error) {
	// get a list of the files in /proc
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Print("Error reading from /proc", err)
		return nil, err
	}

	// loop over files in /proc; extract info for each one that's a process
	procs := make([]Process, 0, 50)
	for _, file := range files {
		// ignore anything that's not a directory
		if !file.IsDir() {
			continue
		}

		// we only care about processes, which have numeric names
		if !unicode.IsDigit([]rune(file.Name())[0]) {
			continue
		}

		// get this proc's PID
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			log.Print("Error converting proc ID to int (process might have exited)", err)
			continue
		}

		// get this proc's name
		commBytes, err := ioutil.ReadFile(filepath.Join("/proc", file.Name(), "comm"))
		if err != nil {
			log.Print("Error reading proc name (process might have exited)", err)
			continue
		}
		command := strings.TrimSpace(string(commBytes))

		// get the command line args
		argBytes, err := ioutil.ReadFile(filepath.Join("/proc", file.Name(), "cmdline"))
		if err != nil {
			log.Print("Error reading command line args (process might have exited)", err)
			continue
		}
		args := strings.Split(strings.TrimSpace(string(argBytes)), "\x00")

		// get the proc's start time
		statBytes, err := ioutil.ReadFile(filepath.Join("/proc", file.Name(), "stat"))
		if err != nil {
			log.Print("Error reading proc stat (process might have exited)", err)
			continue
		}
		stats := strings.Fields(string(statBytes))
		startTime, err := strconv.ParseUint(stats[startTimeFieldNum], 10, 64)
		if err != nil {
			log.Print("Error parsing proc start time", err)
			continue
		}

		// if everything succeeded for this proc, make a new Process and add to slice
		procs = append(procs, Process{Pid: pid, Command: command, Args: args, StartTime: startTime, status: status})
	}

	return procs, nil
}
