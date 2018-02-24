package main

import (
	"fmt"
	"github.com/patwie/cluster-top/proc"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
	"time"
)

// hash map of all processes
var work_processes map[int]*proc.Process

var clus Cluster

func main() {

	work_processes = make(map[int]*proc.Process)

	// load ports and ip-address
	cfg := LoadConfig()
	cfg.Print()

	// sending messages (PUSH-PULL)
	SocketAddr := "tcp://" + cfg.RouterIp + ":" + cfg.Ports.Nodes
	log.Println("Now pushing to", SocketAddr)
	socket, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	socket.Connect(SocketAddr)

	node := Node{}
	InitNode(&node)
	clus.Nodes = append(clus.Nodes, node)

	cpu_tick_prev := int64(0)
	cpu_tick_cur := int64(0)
	cores := proc.NumberCPUCores()

	fmt.Printf("Found %v cores\n", cores)

	for {
		// reset most processes
		proc.MarkDirtyProcessList(work_processes)

		cpu_tick_cur = proc.CpuTick()
		proc.UpdateProcessList(work_processes)
		clus.Nodes[0].Cpu.Update()

		factor := float32(cpu_tick_cur-cpu_tick_prev) / float32(cores) / 100.
		clus.Nodes[0].Processes = GetProcesses(work_processes, factor, cfg.MaxDisplay)

		FetchMemory(&clus.Nodes[0].Memory)
		clus.Nodes[0].Time = time.Now()

		// encode data
		msg, err := msgpack.Marshal(&clus.Nodes[0])
		if err != nil {
			log.Fatal("encode error:", err)
			panic(err)
		}

		// send data
		socket.SendBytes(msg, 0)

		cpu_tick_prev = cpu_tick_cur
		time.Sleep(time.Duration(cfg.Tick) * time.Second)

	}

}
