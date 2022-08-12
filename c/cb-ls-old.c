#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <string.h>

#define LOTS 2000

struct cb_dir {
        ino_t ino;
        char* name;
    } cb_dirs[LOTS];

// params: directory to open, empty array (pointer to start of array), pointer to current length of array
void getdirs(char *, struct cb_dir *, int *);

int main() {
    int i;
    int *ip;

    i = 0;
    ip = &i;
    getdirs(".", cb_dirs, ip);

    printf("len is %d, and first name is...%s.\n", i, cb_dirs[0].name);

    if (i > LOTS) {
        fprintf(stderr, "More dirs than slots for dirs.");
        return 1;
    }
    return 0;
}

void getdirs(char *dirstr, struct cb_dir *dirs, int *ip) {
    DIR *d;
    struct dirent *l;

    d = opendir(dirstr);

    while((l = readdir(d)) != NULL) {
        struct cb_dir dir = { l->d_fileno, l->d_name };
        *dirs = dir;
        dirs++;
        (*ip)++;
    }
}
