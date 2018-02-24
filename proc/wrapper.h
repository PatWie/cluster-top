// Author: Patrick Wieschollek, 2018
#ifndef GOCODE_WRAPPER_H
#define GOCODE_WRAPPER_H

void clock_ticks(long int *hz);
void time_of_day(float *current_time);
int boot_time(float *uptime, float *idle);
void get_mem(unsigned long *mem_total, unsigned long *mem_free, unsigned long *mem_available);
unsigned long long int read_cpu_tick();
void get_uid_from_pid(unsigned long pid, unsigned long *uid);
void read_pid_info(unsigned long pid, unsigned long *time, unsigned long long *starttime, char *name);
unsigned int num_cores();
void read_cpu_info(unsigned long long int* total_time, unsigned long long int* ioWait);
#endif // GOCODE_WRAPPER_H