package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockStatusGood1(pid int) (map[string]string, error) {
	return map[string]string{
		"CapInh":                     "0000000000000000",
		"Seccomp":                    "0",
		"State":                      "S (sleeping)",
		"Tgid":                       "24263",
		"Ngid":                       "0",
		"Uid":                        "0\t0\t0\t0",
		"Threads":                    "3",
		"SigCgt":                     "00000081800006e8",
		"Mems_allowed":               "00000000,00000001",
		"Mems_allowed_list":          "0",
		"VmPeak":                     "2309960 kB",
		"VmSize":                     "1353584 kB",
		"VmLck":                      "0 kB",
		"VmExe":                      "12340 kB",
		"VmLib":                      "68656 kB",
		"Cpus_allowed":               "2",
		"Gid":                        "0\t0\t0\t0",
		"NSpgid":                     "24263",
		"VmPTE":                      "444 kB",
		"HugetlbPages":               "1048576 kB",
		"SigQ":                       "0/15649",
		"CapBnd":                     "0000003fffffffff",
		"FDSize":                     "64",
		"NSpid":                      "24263",
		"VmRSS":                      "54584 kB",
		"voluntary_ctxt_switches":    "3",
		"nonvoluntary_ctxt_switches": "13",
		"CapAmb":                     "0000000000000000",
		"TracerPid":                  "0",
		"Groups":                     "0",
		"NSsid":                      "24263",
		"VmPin":                      "0 kB",
		"ShdPnd":                     "0000000000000000",
		"SigBlk":                     "0000000000000000",
		"PPid":                       "1",
		"NStgid":                     "24263",
		"VmStk":                      "2180 kB",
		"CapEff":                     "0000003fffffffff",
		"Name":                       "bessd",
		"VmData":                     "104832 kB",
		"CapPrm":                     "0000003fffffffff",
		"Cpus_allowed_list":          "1",
		"Pid":                        "24263",
		"VmHWM":                      "54584 kB",
		"VmPMD":                      "20 kB",
		"VmSwap":                     "0 kB",
		"SigPnd":                     "0000000000000000",
		"SigIgn":                     "0000000001001000"}, nil
}

func mockStatusGood2(pid int) (map[string]string, error) {
	return map[string]string{
		"NSsid":  "1887",
		"VmExe":  "3836 kB",
		"VmLib":  "2040 kB",
		"SigPnd": "0000000000000000",
		"State":  "S (sleeping)",
		"Pid":    "26690",
		"VmHWM":  "11552 kB",
		"nonvoluntary_ctxt_switches": "5",
		"SigQ":                    "0/15649",
		"CapPrm":                  "0000000000000000",
		"CapBnd":                  "0000003fffffffff",
		"Tgid":                    "26690",
		"TracerPid":               "0",
		"VmLck":                   "0 kB",
		"VmRSS":                   "11552 kB",
		"VmPMD":                   "20 kB",
		"ShdPnd":                  "0000000000000000",
		"SigIgn":                  "0000000000000000",
		"Cpus_allowed_list":       "0-1",
		"Name":                    "go",
		"VmSize":                  "72332 kB",
		"VmPTE":                   "88 kB",
		"HugetlbPages":            "0 kB",
		"Threads":                 "7",
		"Mems_allowed_list":       "0",
		"FDSize":                  "64",
		"NStgid":                  "26690",
		"VmStk":                   "2180 kB",
		"Seccomp":                 "0",
		"Ngid":                    "0",
		"Groups":                  "4 24 27 30 46 110 115 116 998 1000",
		"VmPin":                   "0 kB",
		"SigBlk":                  "0000000000000000",
		"voluntary_ctxt_switches": "201",
		"SigCgt":                  "ffffffffffc1feff",
		"CapInh":                  "0000000000000000",
		"Cpus_allowed":            "3",
		"PPid":                    "26689",
		"Uid":                     "1000\t1000\t1000\t1000",
		"Gid":                     "1000\t1000\t1000\t1000",
		"NSpgid":                  "4417",
		"VmPeak":                  "72332 kB",
		"Mems_allowed":            "00000000,00000001",
		"NSpid":                   "26690",
		"VmData":                  "57240 kB",
		"VmSwap":                  "0 kB",
		"CapEff":                  "0000000000000000",
		"CapAmb":                  "0000000000000000"}, nil
}

func mockStatusMemUnitsBytes(pid int) (map[string]string, error) {
	return map[string]string{
		"VmRSS":        "11552 B",
		"VmSize":       "72332 B",
		"HugetlbPages": "0 B"}, nil
}

func mockStatusMemMixedUnits(pid int) (map[string]string, error) {
	return map[string]string{
		"VmRSS":        "11552 B",
		"VmSize":       "72332 kB",
		"HugetlbPages": "123 MB"}, nil
}

func mockStatusMemMissingHuge(pid int) (map[string]string, error) {
	return map[string]string{
		"VmRSS":  "11552 kB",
		"VmSize": "72332 kB"}, nil
}

func mockStatusMemMalformed(pid int) (map[string]string, error) {
	return map[string]string{
		"VmRSS":        "xyz kB",
		"VmSize":       "123 456",
		"HugetlbPages": "0kB"}, nil
}

func mockStatusStateNoSpace(pid int) (map[string]string, error) {
	return map[string]string{
		"State": "S(sleeping)"}, nil
}

func mockStatusStateTab(pid int) (map[string]string, error) {
	return map[string]string{
		"State": "S\t(sleeping)"}, nil
}

func mockStatusStateNoParens(pid int) (map[string]string, error) {
	return map[string]string{
		"State": "S sleeping"}, nil
}

func mockStatusStateMissing(pid int) (map[string]string, error) {
	return map[string]string{
		"VmRSS":        "11552 B",
		"VmSize":       "72332 kB",
		"HugetlbPages": "123 MB"}, nil
}

func TestParseMemorySize(t *testing.T) {
	size, err := parseMemorySize("1 kB")
	assert.Nil(t, err)
	assert.Equal(t, uint64(1024), size)

	size, err = parseMemorySize("1 B")
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), size)

	size, err = parseMemorySize(" 1 kB ")
	assert.Nil(t, err)
	assert.Equal(t, uint64(1024), size)

	size, err = parseMemorySize("")
	assert.EqualError(t, err, "memory size string not of form '1353584 kB'")

	size, err = parseMemorySize("123 ")
	assert.EqualError(t, err, "memory size string not of form '1353584 kB'")

	size, err = parseMemorySize("VmSize:  1353584 kB")
	assert.EqualError(t, err, "memory size string not of form '1353584 kB'")

	size, err = parseMemorySize("xyz kB")
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"xyz\": invalid syntax")

	size, err = parseMemorySize("-1 kB")
	assert.EqualError(t, err, "strconv.ParseUint: parsing \"-1\": invalid syntax")

	size, err = parseMemorySize("1 xyz")
	assert.EqualError(t, err, "unrecognized units: xyz")
}

func TestMemUsage(t *testing.T) {

	// test with mocked processes
	pGood1 := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusGood1}
	vmSizeBytes, rssSizeBytes, hugeSizeBytes, err := pGood1.MemUsage()
	assert.Nil(t, err)
	assert.Equal(t, uint64(1353584*1024), vmSizeBytes)
	assert.Equal(t, uint64(54584*1024), rssSizeBytes)
	assert.Equal(t, uint64(1048576*1024), hugeSizeBytes)

	pGood2 := Process{Pid: 26690, Command: "", Args: nil, StartTime: 1000, status: mockStatusGood2}
	vmSizeBytes, rssSizeBytes, hugeSizeBytes, err = pGood2.MemUsage()
	assert.Nil(t, err)
	assert.Equal(t, uint64(72332*1024), vmSizeBytes)
	assert.Equal(t, uint64(11552*1024), rssSizeBytes)
	assert.Equal(t, uint64(0), hugeSizeBytes)

	pUnitsBytes := Process{Pid: 26690, Command: "", Args: nil, StartTime: 1000, status: mockStatusMemUnitsBytes}
	vmSizeBytes, rssSizeBytes, hugeSizeBytes, err = pUnitsBytes.MemUsage()
	assert.Nil(t, err)
	assert.Equal(t, uint64(72332), vmSizeBytes)
	assert.Equal(t, uint64(11552), rssSizeBytes)
	assert.Equal(t, uint64(0), hugeSizeBytes)

	pMissingHuge := Process{Pid: 26690, Command: "", Args: nil, StartTime: 1000, status: mockStatusMemMissingHuge}
	vmSizeBytes, rssSizeBytes, hugeSizeBytes, err = pMissingHuge.MemUsage()
	assert.Nil(t, err)
	assert.Equal(t, uint64(72332*1024), vmSizeBytes)
	assert.Equal(t, uint64(11552*1024), rssSizeBytes)
	assert.Equal(t, uint64(0), hugeSizeBytes)

	pMixedUnits := Process{Pid: 26690, Command: "", Args: nil, StartTime: 1000, status: mockStatusMemMixedUnits}
	vmSizeBytes, rssSizeBytes, hugeSizeBytes, err = pMixedUnits.MemUsage()
	assert.Nil(t, err)
	assert.Equal(t, uint64(72332*1024), vmSizeBytes)
	assert.Equal(t, uint64(11552), rssSizeBytes)
	assert.Equal(t, uint64(0), hugeSizeBytes)

	pMalformed := Process{Pid: 26690, Command: "", Args: nil, StartTime: 1000, status: mockStatusMemMalformed}
	vmSizeBytes, rssSizeBytes, hugeSizeBytes, err = pMalformed.MemUsage()
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), vmSizeBytes)
	assert.Equal(t, uint64(0), rssSizeBytes)
	assert.Equal(t, uint64(0), hugeSizeBytes)

	// test with real processes
	procs, err := Processes()
	assert.Nil(t, err)

	for _, proc := range procs {
		// go should be running, so we'll use it for our tests
		if proc.Command == "go" {
			vmSizeBytes, rssSizeBytes, hugeSizeBytes, err := proc.MemUsage()
			assert.Nil(t, err)
			assert.True(t, vmSizeBytes >= 0, "vmSize should be non-negative")
			assert.True(t, rssSizeBytes >= 0, "rssSize should be non-negative")
			assert.True(t, hugeSizeBytes >= 0, "hugeSize should be non-negative")
		}
	}
}

func TestState(t *testing.T) {
	pGood1 := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusGood1}
	state, description, err := pGood1.State()
	assert.Nil(t, err)
	assert.Equal(t, "S", state)
	assert.Equal(t, "sleeping", description)

	pGood2 := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusGood2}
	state, description, err = pGood2.State()
	assert.Nil(t, err)
	assert.Equal(t, "S", state)
	assert.Equal(t, "sleeping", description)

	pNoSpace := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusStateNoSpace}
	state, description, err = pNoSpace.State()
	assert.EqualError(t, err, "'State' not in expected format: S(sleeping)")

	pTab := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusStateTab}
	state, description, err = pTab.State()
	assert.Nil(t, err)
	assert.Equal(t, "S", state)
	assert.Equal(t, "sleeping", description)

	pNoParens := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusStateNoParens}
	state, description, err = pNoParens.State()
	assert.Nil(t, err)
	assert.Equal(t, "S", state)
	assert.Equal(t, "sleeping", description)

	pMissing := Process{Pid: 24263, Command: "", Args: nil, StartTime: 1000, status: mockStatusStateMissing}
	state, description, err = pMissing.State()
	assert.EqualError(t, err, "'State' not found in process info")
}

func TestProcesses(t *testing.T) {
	procs, err := Processes()

	assert.Nil(t, err)
	assert.NotEmpty(t, procs, "The list of processes should not be empty (*some* procs must be running)")

	foundGoProc := false
	for _, proc := range procs {
		if proc.Command == "go" {
			foundGoProc = true
		}

		assert.True(t, proc.Pid > 0, "process ID should be positive")
		assert.True(t, proc.StartTime > 0, "process start time should be positive")
	}
	assert.True(t, foundGoProc, "go should be one of the running processes")
}

func TestStatus(t *testing.T) {
	procs, err := Processes()
	assert.Nil(t, err)

	for _, proc := range procs {
		// go should be running, so we'll use it for our tests
		if proc.Command == "go" {
			status, err := proc.status(proc.Pid)
			assert.Nil(t, err) // proc should not have exited
			assert.NotEmpty(t, status)

			name, ok := status["Name"]
			assert.True(t, ok, "Name should be a field in /proc/[pid]/status")
			assert.Equal(t, "go", name)
		}
	}
}
