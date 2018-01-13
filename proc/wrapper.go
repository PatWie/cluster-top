// Author: Patrick Wieschollek, 2018
package proc

// #include "wrapper.h"
import "C"

func CpuTick() (t int64) {
	return int64(C.read_cpu_tick())
}

func PidInformation(pid int) (t int64, n string) {
	time := C.ulong(0)
	buf := string("                                                                                                                                ")
	c_dst := C.CString(buf)
	C.read_time_and_name_from_pid(C.ulong(pid), &time, c_dst)

	return int64(time), C.GoString(c_dst)
}

func NumCores() (n int) {
	return int(C.num_cores())
}

func UidFromPid(pid int) (uid int) {
	c_uid := C.ulong(0)
	C.get_uid_from_pid(C.ulong(pid), &c_uid)

	return int(c_uid)
}

func GetMemoryInfo() (total int64, free int64, available int64) {
	c_total := C.ulong(0)
	c_free := C.ulong(0)
	c_available := C.ulong(0)

	C.get_mem(&c_total, &c_free, &c_available)

	return int64(c_total), int64(c_free), int64(c_available)
}
