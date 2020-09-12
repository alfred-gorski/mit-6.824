package mr

import (
	"hash/fnv"
	"log"
	"net/rpc"
)

// KeyValue Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// ihash use ihash(key) % NReduce to choose the reduce task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// Worker called by main/mrworker.go to send an RPC to the master asking for a task
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	// TODO:

	args := Args{}
	reply := Reply{}

	for {
		if success := call("Master.HandlTask", &args, &reply); !success {
			continue
		}
		switch reply.taskType {
		case Map:
		case Reduce:
		}

	}

}

/*
call send an RPC request to the master, wait for the response. Usually returns true, returns false if something goes wrong.
*/
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()
	err = c.Call(rpcname, args, reply)
	if err != nil {
		log.Fatal("calling:", err)
		return false
	}
	return true
}
