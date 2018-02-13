package proc

/*
Patrick Wieschollek, 2018
*/

import (
	"io/ioutil"
	"strconv"
)

// representing a single process
type Process struct {
	UID      int
	TimePrev int64
	Dirty    bool
	Active   bool
	Usage    float32
	PIDInfo  PIDInfo
}

func (p *Process) CurrentTime() int64 {
	return p.PIDInfo.UsedTime + p.PIDInfo.StartTime
}

func UpdateProcessList(procs map[int]*Process) {

	// gather all possible pids
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		// list all possible directories
		name := file.Name()
		// we are just interested in numerical names
		if name[0] < '0' || name[0] > '9' {
			continue
		}
		// get pid
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		p := procs[pid]

		uid := UIDFromPID(pid)

		// ignore root
		if uid == 0 {
			continue
		}

		if p == nil {
			// is a new process
			p = &Process{UID: uid,
				TimePrev: 0,
				Dirty:    true,
				Active:   true,
				Usage:    0.,
				PIDInfo:  PIDInfo{},
			}
		} else {
			// just update
			p.Dirty = false
			p.Active = true
			p.TimePrev = p.CurrentTime()
		}

		p.PIDInfo = InfoFromPid(pid)
		procs[pid] = p
	}

	// remove all processes which are not active anymore
	// pless GO as the following is safe
	for key, v := range procs {
		if v.Active == false {
			delete(procs, key)
		}
	}
}

type ByUsage []Process

func (a ByUsage) Len() int      { return len(a) }
func (a ByUsage) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsage) Less(i, j int) bool {
	return a[i].Usage > a[j].Usage
}

func MarkDirtyProcessList(procs map[int]*Process) {
	for k, _ := range procs {
		procs[k].Dirty = true
		procs[k].Active = false
	}
}
