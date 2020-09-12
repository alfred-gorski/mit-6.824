package mr

// RPC definitions.
// remember to capitalize all names.

import (
	"os"
	"strconv"
)

// Args show how to declare the arguments for an RPC.
type Args struct {
}

// Reply show how to declare reply for an RPC.
type Reply struct {
	taskType    taskType
	index       int
	fileContent string
	fileName    string
}

// TaskType definde
type (
	taskType   uint8
	taskStatus uint8
)

// MAP REDUCE DIE
const (
	Map taskType = iota
	Reduce
	Die
)

// idle task can be rescheduled
const (
	idle taskStatus = iota
	inProgress
	completed
)

/*
Cook up a unique-ish UNIX-domain socket name in /var/tmp, for the master. Can't use the current directory since Athena AFS doesn't support UNIX-domain sockets.
*/
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
