package main

import (
	"github.com/patwie/cluster-top/proc"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
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

// hash map of all processes
var m map[int]*proc.Process

func GetProcesses(procs map[int]*proc.Process, factor float32, max int) []Process {

	var m_display []Process

	for _, v := range procs {
		if v.Dirty == false && v.Active == true {
			if v.TimeCur-v.TimePrev > 0 {
				usage := 1. / factor * float32(v.TimeCur-v.TimePrev)
				copy_proc := Process{v.Pid, v.Name, v.Uid, "", usage}
				m_display = append(m_display, copy_proc)
			}
		}
	}

	sort.Sort(ByUsage(m_display))
	m_display = m_display[:minimum(max, len(m_display))]

	for i := 0; i < len(m_display); i++ {
		user, _ := user.LookupId(strconv.Itoa(m_display[i].Uid))
		m_display[i].Username = user.Username
	}

	return m_display

}

func main() {

	// load ports and ip-address
	cfg := CreateConfig()

	// sending messages (PUSH-PULL)
	SocketAddr := "tcp://" + cfg.ServerIp + ":" + cfg.ServerPortGather
	log.Println("Now pushing to", SocketAddr)
	socket, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	socket.Connect(SocketAddr)

	Machine := &Node{}
	InitNode(Machine)

	m = make(map[int]*proc.Process)

	cpu_tick_prev := int64(0)
	cpu_tick_cur := int64(0)
	cores := proc.NumCores()

	first := true

	for {
		// reset most
		proc.MarkDirtyProcessList(m)

		cpu_tick_cur = proc.CpuTick()
		proc.UpdateProcessList(m)
		factor := float32(cpu_tick_cur-cpu_tick_prev) / float32(cores) / 100.

		Machine.Processes = GetProcesses(m, factor, cfg.MaxDisplay)

		FetchMemory(&Machine.Memory)

		if first != true {
			// provide
			// encode data
			msg, err := msgpack.Marshal(&Machine)
			if err != nil {
				log.Fatal("encode error:", err)
				panic(err)
			}
			// send data
			socket.SendBytes(msg, 0)
		}

		first = false

		cpu_tick_prev = cpu_tick_cur
		time.Sleep(cfg.Tick)

	}

}
