/*
Package mr impliments MapReduce worker, master, rpc
*/
package mr

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

// Master defines data structure of master
type Master struct {
	inputFiles  []string
	mapTasks    *Tasks
	reduceTasks *Tasks
}

// Tasks defines Task arguments for Master. num is the number of Tasks,
type Tasks struct {
	num   int
	mutex sync.Mutex
	task  []Task
}

type Task struct {
	stat taskStatus
}

func readFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	return string(content)

}

func (t *Tasks) finished() bool {
	for i := 0; i < t.num; i++ {
		if t.task[i].stat != completed {
			return false
		}
	}
	return true
}

// HandleTask ask for a Task from the master
func (m *Master) HandleTask(args *Args, reply *Reply) error {
	//TODO: RPC handlers for the worker to call.
	m.mapTasks.mutex.Lock()
	defer m.mapTasks.mutex.Unlock()

	for i := 0; i < m.mapTasks.num; i++ {
		if m.mapTasks.task[i].stat == idle {
			reply = &Reply{
				taskType:    Map,
				index:       i,
				fileContent: readFile(m.inputFiles[i]),
			}

		}
	}

	if m.mapTasks.finished() {
		
	}

	return nil
}

// start a thread that listens for RPCs from worker.go
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// Done main/mrmaster.go calls Done() periodically to find out if the entire job has finished.
func (m *Master) Done() bool {

	// TODO: Your code here.

	// Map finished

	// Reduce finished

	return true
}

// MakeMaster create a Master. main/mrmaster.go calls this function. nReduce is the number of reduce tasks to use.
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{
		inputFiles:  files,
		mapTasks:    NewTasks(len(files)),
		reduceTasks: NewTasks(nReduce),
	}

	m.server()
	return &m
}

// NewTask create a task
func NewTasks(n int) *Tasks {
	tasks := Tasks{num: n}
	return &tasks
}
