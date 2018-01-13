#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <unistd.h>



void get_uid_from_pid(long tgid, unsigned long *uid) {
    char path[40], line[100], *p;
    FILE* statusf;

    snprintf(path, 40, "/proc/%ld/status", tgid);

    statusf = fopen(path, "r");
    if(!statusf)
        return;

    while(fgets(line, 100, statusf)) {
        if(strncmp(line, "Uid:", 4) != 0)
            continue;
        // Uid: 1000    1000    1000    1000
        sscanf(line, "%*s %lu %*s", uid);
        break;
    }
    fclose(statusf);
}


int main(int argc, char const *argv[])
{
    int pid = 29575;
    unsigned long uid = 0;
    get_uid_from_pid(pid, &uid);
    printf("uid: %lu\n", uid);
    return 0;
}