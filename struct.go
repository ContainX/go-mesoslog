package main

// LogType - enumeration of STDOUT or STDERR
type LogType int

const (
	// STDOUT - requests the stdout log from Mesos
	STDOUT LogType = 1 + iota
	// STDERR - requests the stderr log from Mesos
	STDERR
)

var logTypes = [...]string{
	"stdout",
	"stderr",
}

func (l LogType) String() string {
	return logTypes[l-1]
}

type masterState struct {
	Frameworks []*mstateFramework `json:"frameworks"`
	Slaves     []*mstateSlave     `json:"slaves"`
}

type mstateFramework struct {
	Tasks []*mstateTask `json:"tasks"`
}

type mstateSlave struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Pid      string `json:"pid"`
}

type mstateTask struct {
	ID          string `json:"id"`
	FrameworkID string `json:"framework_id"`
	ExecutorID  string `json:"executor_id"`
	SlaveID     string `json:"slave_id"`
	Name        string `json:"name"`
}

type slaveState struct {
	ID         string             `json:",id"`
	Frameworks []*sstateFramework `json:",frameworks"`
}

type sstateFramework struct {
	ID        string            `json:"id"`
	Executors []*sstateExecutor `json:"executors"`
}

type sstateExecutor struct {
	ID        string `json:"id"`
	Directory string `json:"directory"`
}

// LogOut - struct which holds the result from getting Mesos logs
//
type LogOut struct {
	// AppID - the task name / application identifier
	AppID string
	// TaskID - the task identifier
	TaskID string
	// Log - filename of the outputted log when in download more or RAW log if request is to print to stdout
	Log string
}
