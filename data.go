package main

import (
	"fmt"
	"github.com/apcera/termtables"
	"github.com/patwie/cluster-top/proc"
	"os"
	"os/user"
	"sort"
	"strconv"
	"time"
)

func minimum(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

type Cluster struct {
	Nodes []Node `json:"nodes"`
}

type Process struct {
	UID      int
	Username string
	Usage    float32
	Info     proc.PIDInfo
}

type Memory struct {
	Total     int64
	Free      int64
	Available int64
	Used      int64
	Usage     float32
}

func FetchMemory(m *Memory) {
	m.Total, m.Free, m.Available = proc.GetRAMMemoryInfo()
	m.Used = m.Total - m.Free
	m.Usage = float32(100 * float64(m.Used) / float64(m.Total))
}

type Cpu struct {
	Cores int
}

type Node struct {
	Name      string
	Processes []Process
	Memory    Memory
	Time      time.Time `json:"time"` // current timestamp from message
	Cpu       Cpu
}

func InitNode(n *Node) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	n.Name = name
	n.Cpu.Cores = proc.NumberCPUCores()
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

func (c *Cluster) Print(show_time bool) {

	table := termtables.CreateTable()

	tableHeader := []interface{}{"Node", "RAM-Util", "RAM-Util", "PID", "User", "Command", "CPU-Util"}
	if show_time {
		tableHeader = append(tableHeader, "Last Seen")
	}
	table.AddHeaders(tableHeader...)

	for n_id, n := range c.Nodes {

		node_lastseen := n.Time.Format("Mon Jan 2 15:04:05 2006")

		if len(n.Processes) == 0 {

			memory := strconv.FormatInt(n.Memory.Used/1024/1024, 10) +
				"GiB / " +
				strconv.FormatInt(n.Memory.Total/1024/1024, 10) + "GiB"
			memory_per := strconv.Itoa(int(n.Memory.Usage)) + "%"

			tableRow := []interface{}{
				n.Name,
				memory,
				memory_per,
				"",
				"",
				"",
				"",
			}

			if show_time {
				tableRow = append(tableRow, node_lastseen)
			}

			table.AddRow(tableRow...)
			table.SetAlign(termtables.AlignRight, 2)
		} else {
			for p_id, p := range n.Processes {
				name := ""
				memory := ""
				memory_per := ""

				if p_id == 0 {
					name = n.Name
					memory = strconv.FormatInt(n.Memory.Used/1024/1024, 10) +
						"GiB / " +
						strconv.FormatInt(n.Memory.Total/1024/1024, 10) + "GiB"
					memory_per = strconv.Itoa(int(n.Memory.Usage)) + "%"
				}

				tableRow := []interface{}{
					name,
					memory,
					memory_per,
					p.Info.PID,
					p.Username,
					p.Info.Command,
					strconv.Itoa(int(p.Usage)) + "%",
				}

				if show_time {
					if p_id == 0 {
						tableRow = append(tableRow, node_lastseen)

					} else {
						tableRow = append(tableRow, "")

					}

				}

				table.AddRow(tableRow...)
				table.SetAlign(termtables.AlignRight, 2)

			}
		}
		if n_id < len(c.Nodes)-1 {
			table.AddSeparator()
		}

	}
	fmt.Printf("\033[2J")
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 2006") + " (http://github.com/patwie/cluster-smi)")
	fmt.Println(table.Render())
}

func GetProcesses(procs map[int]*proc.Process, factor float32, max int) []Process {

	var m_display []Process

	for _, v := range procs {
		if v.Dirty == false && v.Active == true {
			if v.CurrentTime()-v.TimePrev > 0 {
				usage := 1. / factor * float32(v.CurrentTime()-v.TimePrev)

				copy_proc := Process{v.UID, "", usage, v.PIDInfo}
				// fmt.Println(copy_proc)
				m_display = append(m_display, copy_proc)
			}
		}
	}

	sort.Sort(ByUsage(m_display))
	m_display = m_display[:minimum(max, len(m_display))]

	for i := 0; i < len(m_display); i++ {
		user, _ := user.LookupId(strconv.Itoa(m_display[i].UID))
		m_display[i].Username = user.Username
	}
	return m_display

}
