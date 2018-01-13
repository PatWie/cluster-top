#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <unistd.h>



void get_mem(unsigned long *mem_total, unsigned long *mem_free, unsigned long *mem_available) {
    char line[100], *p;
    FILE* statusf;


    statusf = fopen("/proc/meminfo", "r");
    if (!statusf)
        return;


    fgets(line, 100, statusf);
    sscanf(line, "%*s %lu %*s", mem_total);
    fgets(line, 100, statusf);
    sscanf(line, "%*s %lu %*s", mem_free);
    fgets(line, 100, statusf);
    sscanf(line, "%*s %lu %*s", mem_available);


    fclose(statusf);
}


int main(int argc, char const *argv[]) {
    unsigned long mem_total = 0;
    unsigned long mem_free = 0;
    unsigned long mem_available = 0;
    get_mem(&mem_total, &mem_free, &mem_available);
    printf("mem_total: %lu\n", mem_total);
    printf("mem_free: %lu\n", mem_free);
    printf("mem_available: %lu\n", mem_available);
    return 0;
}