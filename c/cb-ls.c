#include <dirent.h>    /* ino_t */
#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>    /* qsort */
#include <unistd.h>    /* getopt */

#define LOTS 2000

struct cb_dir {
        ino_t ino;
        char* name;
    } cb_dirs[LOTS];


int dir_comp(const void *v1, const void *v2) {
    const struct cb_dir *d1 = v1;
    const struct cb_dir *d2 = v2;

    // Returns <1, 0, or >1 after comparing the strings
    // Note the use of -> to get the name property off of the struct
    // Remember that we're using pointers to our array elements, so
    // to get the properties from those elements, we could deference and
    // then use the . to get the value, or just use ->
    // like return strcmp((*d1).name, (*d2).name);
    return strcmp(d1->name, d2->name);
}
// params: directory to open, empty array (pointer to start of array), pointer to current length of array
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

int main(int argc, char *argv[]) {
    int i, opt, all, list;
    int *ip;
    char *dir;

    i = 0;
    ip = &i;
    // TODO add help
    while ((opt = getopt(argc, argv, "al")) != -1) {
        switch(opt) {
        case 'a':
            all = 1;
            break;
        case 'l':
            list = 1;
            break;
        }
    }
    printf("show all? %d\nshow list? %d\n", all, list);

    if (argc == 2) {
        dir = argv[1];
    } else {
        dir = ".";
    }

    getdirs(dir, cb_dirs, ip);

    qsort(cb_dirs, i, sizeof(struct cb_dir), dir_comp);

    for(int j = 0; j < i; j++) {
        printf("name is %s\n", cb_dirs[j].name);
    }

    if (i > LOTS) {
        fprintf(stderr, "More dirs than slots for dirs.");
        return 1;
    }
    return 0;
}

