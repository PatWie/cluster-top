package main

import (
	"flag"
	"github.com/patwie/cluster-top/proc"
	"time"
)

// hash map of all processes
var work_processes map[int]*proc.Process

var clus Cluster

func main() {

	showTimePtr := flag.Bool("t", false, "show time of events")
	flag.Parse()

	work_processes = make(map[int]*proc.Process)

	// load ports and ip-address
	cfg := LoadConfig()

	node := Node{}
	InitNode(&node)
	clus.Nodes = append(clus.Nodes, node)

	cpu_tick_prev := int64(0)
	cpu_tick_cur := int64(0)
	cores := proc.NumberCPUCores()

	for {
		// reset most processes
		proc.MarkDirtyProcessList(work_processes)

		cpu_tick_cur = proc.CpuTick()
		proc.UpdateProcessList(work_processes)

		factor := float32(cpu_tick_cur-cpu_tick_prev) / float32(cores) / 100.
		clus.Nodes[0].Processes = GetProcesses(work_processes, factor, cfg.MaxDisplay)

		FetchMemory(&clus.Nodes[0].Memory)
		clus.Nodes[0].Time = time.Now()

		clus.Print(*showTimePtr)

		cpu_tick_prev = cpu_tick_cur
		time.Sleep(time.Duration(cfg.Tick) * time.Second)

	}

}
