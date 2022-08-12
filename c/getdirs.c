#include <dirent.h>
#include <stdio.h>
#define LOTS 2000

struct cb_dir {
        ino_t ino;
        char* name;
    } cb_dirs[LOTS];

// params: directory to open, empty array (pointer to start of array), current working index of array (init to zero)
int getdirs(char *, struct cb_dir *, int);

int main() {
    int i, len;

    i = 0;
    len = getdirs(".", cb_dirs, i);
    printf("len is %d, and first name is...%s.\n", len, cb_dirs[8].name);
    if (len > LOTS) {
        fprintf(stderr, "More dirs than slots for dirs.");
        return 1;
    }
    return 0;
}

int getdirs(char *dirstr, struct cb_dir *dirs, int i) {
    DIR *d;
    struct dirent *l;

    d = opendir(dirstr);

    while((l = readdir(d)) != NULL) {
        struct cb_dir dir = { l->d_fileno, l->d_name };
        cb_dirs[i] = dir; // How to do with pointer?
        ++i;
    }
    return i;
}
