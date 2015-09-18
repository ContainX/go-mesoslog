package main

type LogType int

const (
	STDOUT LogType = 1 + iota
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
	Frameworks []*mstateFramework `json:",frameworks"`
	Slaves     []*mstateSlave     `json:",slaves"`
}

type mstateFramework struct {
	Tasks []*mstateTask `json:",tasks"`
}

type mstateSlave struct {
	ID       string `json:",id"`
	Hostname string `json:",hostname"`
	Pid      string `json:",pid"`
}

type mstateTask struct {
	ID           string `json:",id"`
	Framework_ID string `json:",framework_id"`
	Executor_ID  string `json:",executor_id"`
	Slave_ID     string `json:",slave_id"`
	Name         string `json:",name"`
}

type slaveState struct {
	ID         string             `json:",id"`
	Frameworks []*sstateFramework `json:",frameworks"`
}

type sstateFramework struct {
	ID        string            `json:",id"`
	Executors []*sstateExecutor `json:",executors"`
}

type sstateExecutor struct {
	ID        string `json:",id"`
	Directory string `json:",directory"`
}

type LogOut struct {
	AppID  string
	TaskID string
	Log    string
}
