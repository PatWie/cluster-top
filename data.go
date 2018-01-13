package main

import (
	"fmt"
	"github.com/apcera/termtables"
	"github.com/patwie/cluster-top/proc"
	"os"
	"sort"
	"strconv"
)

type Cluster struct {
	Nodes []Node `json:"nodes"`
}

type Process struct {
	Pid      int
	Name     string
	Uid      int
	Username string
	Usage    float32
}

type Memory struct {
	Total     int64
	Free      int64
	Available int64
	Used      int64
	Usage     float32
}

func FetchMemory(m *Memory) {
	m.Total, m.Free, m.Available = proc.GetMemoryInfo()
	m.Used = m.Total - m.Free
	m.Usage = float32(100 * float64(m.Used) / float64(m.Free))
}

type Cpu struct {
	Cores int
}

type Node struct {
	Name      string
	Processes []Process
	Memory    Memory
	Cpu       Cpu
}

func InitNode(n *Node) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	n.Name = name
	n.Cpu.Cores = proc.NumCores()
}

type ByUsage []Process

func (a ByUsage) Len() int      { return len(a) }
func (a ByUsage) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsage) Less(i, j int) bool {
	return a[i].Usage > a[j].Usage
}

type ByName []Node

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (c *Cluster) Sort() {
	sort.Sort(ByName(c.Nodes))
}

func (c *Cluster) Print() {
	table := termtables.CreateTable()

	table.AddHeaders("Node", "RAM-Usage", "Pid", "User", "Command", "CPU-Util")
	for n_id, n := range c.Nodes {
		for p_id, p := range n.Processes {
			name := ""
			memory := ""
			if p_id == 0 {
				name = n.Name
				memory = strconv.FormatInt(n.Memory.Used/1024, 10) +
					"MiB / " +
					strconv.FormatInt(n.Memory.Total/1024, 10) + "MiB" + " (" +
					strconv.Itoa(int(100-n.Memory.Usage)) + "%)"
			}

			table.AddRow(
				name,
				memory,
				p.Pid,
				p.Username,
				p.Name,
				strconv.Itoa(int(p.Usage))+"%",
			)
			if n_id < len(c.Nodes)-1 {
				table.AddSeparator()
			}
		}
	}
	fmt.Printf("\033[2J")
	fmt.Println(table.Render())
}
