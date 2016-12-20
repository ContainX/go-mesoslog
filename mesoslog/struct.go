package mesoslog

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

type taskInfo struct {
	Tasks          []*mstateTask
	CompletedTasks []*mstateTask
}

type mstateFramework struct {
	Tasks          []*mstateTask `json:"tasks"`
	CompletedTasks []*mstateTask `json:"completed_tasks"`
}

type mstateSlave struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Pid      string `json:"pid"`
}

type mstateTask struct {
	ID            string              `json:"id"`
	FrameworkID   string              `json:"framework_id"`
	ExecutorID    string              `json:"executor_id"`
	SlaveID       string              `json:"slave_id"`
	Name          string              `json:"name"`
	Statuses      []*mstateTaskStatus `json:"statuses"`
	LastTimestamp float64             `json:"-"`
	LastState     string              `json:"-"`
}

type mstateTaskStatus struct {
	State     string  `json:"state"`
	Timestamp float64 `json:"timestamp"`
}

type slaveState struct {
	ID         string             `json:",id"`
	Frameworks []*sstateFramework `json:",frameworks"`
}

type sstateFramework struct {
	ID                 string            `json:"id"`
	Executors          []*sstateExecutor `json:"executors"`
	CompletedExecutors []*sstateExecutor `json:"completed_executors"`
}

type sstateExecutor struct {
	ID        string `json:"id"`
	Directory string `json:"directory"`
}

type slaveInfo struct {
	Slave     *mstateSlave
	State     *slaveState
	Directory string
}

type readData struct {
	Data   string `json:"data"`
	Offset int    `json:"offset"`
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

func (t *mstateTask) UpdateLastState(state *mstateTaskStatus) {
	if state != nil {
		t.LastState = state.State
		t.LastTimestamp = state.Timestamp
	}
}
