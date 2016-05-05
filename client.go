package main

// Package sandboxfinder helps to look up the sandbox page on the mesos-master
// where the stdout and stderr streams for a running task can be viewed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// PageLength is the amount of data we want to consume during log tailing
	PageLength int    = 5000
	TailURIFmt string = "http://%s:5051/files/read.json?path=%s&offset=%v&length=%v"
)

// MesosClient holds state about the current Master node.  Allows method receivers to obtain these values
type MesosClient struct {
	Host      string
	Port      int
	MasterURL string
	State     *masterState
}

// NewMesosClient - Creates a new MesosClient
// {host} - the host/ip of the mesos master node
// {port} - the port # of the mesos master node
func NewMesosClient(host string, port int) (*MesosClient, error) {
	masterURL, err := getMasterRedirect(host, port)
	if err != nil {
		return nil, err
	}

	state, err := getMasterState(masterURL)
	if err != nil {
		return nil, err
	}
	return &MesosClient{Host: host, Port: port, MasterURL: masterURL, State: state}, nil
}

// GetAppNames - List all unique app names aka task names running in the Mesos cluster
func (c *MesosClient) GetAppNames() (map[string]int, error) {
	apps := findApps(c.State)
	if apps == nil || len(apps) == 0 {
		return nil, fmt.Errorf("Applications could not be retrieved")
	}
	return apps, nil
}

// GetLog - Gets/Downloads logs for a [appID]
// {appID} - the task name / app identifier
// {logtype} - the desired log type STDOUT | STDERR
// {dir} - optional output dir which is used to download vs stdout
func (c *MesosClient) GetLog(appID string, logtype LogType, dir string) ([]*LogOut, error) {
	var result []*LogOut
	tasks := findTask(c.State, appID)
	if tasks == nil || len(tasks) == 0 {
		return nil, fmt.Errorf("application could not be found")
	}

	for _, task := range tasks {

		slaveInfo, err := c.getSlaveInfo(task)
		if err != nil {
			return nil, err
		}

		url := fmt.Sprintf("http://%s:5051/files/download.json?path=%s/", slaveInfo.Slave.Hostname, slaveInfo.Directory)

		var filename string
		if dir != "" {
			filename = filepath.Join(dir, fmt.Sprintf("%s_%s.txt", task.ID, logtype.String()))
		}
		data, err := download(url+logtype.String(), filename)
		if err != nil {
			return nil, err
		}
		result = append(result, &LogOut{TaskID: task.ID, AppID: appID, Log: data})
	}
	return result, nil
}

// TailLog - Tails the logs for a [appID]
// {appID} - the task name / app identifier
// {logtype} - the desired log type STDOUT | STDERR
// {duration} - poll frequency in seconds
func (c *MesosClient) TailLog(appID string, logtype LogType, duration int) error {
	tasks := findTask(c.State, appID)
	if tasks == nil || len(tasks) == 0 {
		return fmt.Errorf("application could not be found")
	}

	var chans []<-chan string
	for _, task := range tasks {

		slaveInfo, err := c.getSlaveInfo(task)
		if err != nil {
			return err
		}
		o := c.asyncTail(task, slaveInfo, logtype, duration)
		chans = append(chans, o)

	}
	output := merge(chans...)
	for {
		fmt.Print(<-output)
	}
	return nil
}

func (c *MesosClient) asyncTail(task *mstateTask, s *slaveInfo, lt LogType, duration int) <-chan string {
	ch := make(chan string)
	path := fmt.Sprintf("%s/%s", s.Directory, lt.String())
	go func() {
		offset := 0
		for {
			url := fmt.Sprintf(TailURIFmt, s.Slave.Hostname, path, offset, PageLength)

			resp, err := download(url, "")
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				continue
			}
			var rd readData
			json.Unmarshal([]byte(resp), &rd)

			if len(rd.Data) < 5 {
				time.Sleep(time.Duration(duration) * time.Second)
				continue
			}
			offset += len(rd.Data)
			ch <- decorateLog(task.ID, rd.Data)
		}
	}()
	return ch
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func decorateLog(name, data string) string {
	lines := strings.Split(data, "\n")
	buf := new(bytes.Buffer)
	taskIdx := strings.Index(name, ".")
	dec := name[0:taskIdx]
	taskId := name[taskIdx+1:]
	dec = fmt.Sprintf("%s.%s", dec, taskId[0:strings.Index(taskId, "-")])
	for _, l := range lines {
		if len(l) > 0 {
			buf.WriteString(fmt.Sprintf("[%s] %s\n", dec, l))
		}
	}
	return buf.String()
}

func (c *MesosClient) getSlaveInfo(task *mstateTask) (*slaveInfo, error) {
	slave := findSlave(c.State, task.SlaveID)
	if slave == nil {
		return nil, fmt.Errorf("invalid state.json; referenced slave not present")
	}

	slaveURL, err := constructSlaveURL(slave)
	if err != nil {
		return nil, err
	}

	slaveState, err := getSlaveState(slaveURL)
	if err != nil {
		return nil, err
	}

	directory := findDirectory(slaveState, task.FrameworkID, task.ID, task.ExecutorID)
	if directory == "" {
		return nil, fmt.Errorf("couldn't locate directory on slave")
	}
	return &slaveInfo{Slave: slave, State: slaveState, Directory: directory}, nil
}

func getMasterRedirect(host string, port int) (string, error) {
	url := fmt.Sprintf("http://%s:%d/master/redirect", host, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	loc, err := resp.Location()
	if err != nil {
		return "", err
	}
	return loc.String(), nil
}

func getMasterState(masterURL string) (*masterState, error) {
	url := masterURL + "/state.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var mstate masterState
	err = json.Unmarshal(data, &mstate)
	if err != nil {
		return nil, err
	}
	return &mstate, nil
}

func findTask(state *masterState, appID string) map[string]*mstateTask {
	m := make(map[string]*mstateTask)
	for _, framework := range state.Frameworks {
		for _, task := range framework.Tasks {
			if task.Name == appID || task.ID == appID {
				m[task.ID] = task
			}
		}
	}
	return m
}

func findApps(state *masterState) map[string]int {
	m := make(map[string]int)
	for _, framework := range state.Frameworks {
		for _, task := range framework.Tasks {
			m[task.Name]++
		}
	}
	return m
}

func findSlave(state *masterState, slaveID string) *mstateSlave {
	for _, slave := range state.Slaves {
		if slave.ID == slaveID {
			return slave
		}
	}
	return nil
}

func constructSlaveURL(slave *mstateSlave) (*url.URL, error) {
	parts := strings.SplitN(slave.Pid, "@", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid slave pid %s", slave.Pid)
	}
	slaveName := parts[0]

	hostAndPort := strings.Split(parts[1], ":")
	port := "80"
	if len(hostAndPort) > 1 {
		port = hostAndPort[1]
	}

	host := fmt.Sprintf("%s:%s", slave.Hostname, port)
	path := fmt.Sprintf("%s/state.json", slaveName)

	return &url.URL{
		Scheme: "http",
		Host:   host,
		Opaque: fmt.Sprintf("//%s/%s", host, path),
	}, nil
}

func getSlaveState(slaveURL *url.URL) (*slaveState, error) {
	req, err := http.NewRequest("GET", slaveURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL = slaveURL
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var sstate slaveState
	err = json.Unmarshal(data, &sstate)
	if err != nil {
		return nil, err
	}
	return &sstate, nil
}

func findDirectory(sstate *slaveState, frameworkID, taskID, executorID string) string {
	for _, framework := range sstate.Frameworks {
		if framework.ID != frameworkID {
			continue
		}
		for _, executor := range framework.Executors {

			if executor.ID == executorID || executor.ID == taskID {
				return executor.Directory
			}
		}
	}
	return ""
}

func download(slaveURL string, filename string) (string, error) {
	resp, err := http.Get(slaveURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if filename != "" {
		if e := writeFile(filename, resp.Body); e != nil {
			return "", e
		}
		return filename, nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

func writeFile(filename string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil

}
