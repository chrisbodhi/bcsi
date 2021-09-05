#include <dirent.h>    /* ino_t */
#include <errno.h>     /* stderr */
#include <stdio.h>
#include <string.h>
#include <stdlib.h>    /* qsort */
#include <unistd.h>    /* getopt */

#define LOTS 2000

const int true = 1;
const int false = 0;

struct cb_dir {
    ino_t ino;
    char* name;
    int hidden;
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
        struct cb_dir dir = { l->d_fileno, l->d_name, false };
        *dirs = dir;
        dirs++;
        (*ip)++;
    }
}

// mark dirs that start with . as hidden
void hidehidden(struct cb_dir *dirs) {
    for (; dirs->name != NULL; dirs++) {
        if (strncmp(dirs->name, ".", 1) == 0) {
            dirs->hidden = true;
        }
    }
}

void printnames(struct cb_dir *dirs) {
    for(; dirs->name != NULL; dirs++) {
        if (!dirs->hidden) {
            printf("%s\t", dirs->name);
        }
    }
}

void printdetails(struct cb_dir *dirs) {
    for (; dirs->name != NULL; dirs++) {
        if (!dirs->hidden) {
            printf("%llu %s\n", dirs->ino, dirs->name);
        }
    }
}

int main(int argc, char *argv[]) {
    int i, opt, all, inodes;
    int *ip;
    char *dir;

    all = false;
    inodes = false;
    i = 0;
    ip = &i;

    while ((opt = getopt(argc, argv, "ahi")) != -1) {
        switch(opt) {
        case 'a':
            all = true;
            break;
        case 'h':
            printf("cb-ls\nA subset of ls\n\nFlags:\n\t-a:\tShow all files\n\t-i:\tShow inodes of files\n\t-----\n\t-h:\tShow this help\n");
            return 0;
        case 'i':
            inodes = true;
            break;
        default:
            // Exit code 1 if illegal option
            return 1;
        }
    }

    if (argc == 3) {
        dir = argv[2];
    } else if (argc == 2 && strncmp(argv[1], "-", 1) != 0) {
        dir = argv[1];
    } else {
        dir = ".";
    }

    // Get dirs
    getdirs(dir, cb_dirs, ip);

    if (i > LOTS) {
        fprintf(stderr, "More dirs than slots for dirs.");
        return 1;
    }

    // Sort dir names
    qsort(cb_dirs, i, sizeof(struct cb_dir), dir_comp);

    // Act on flags
    // if *not* showing all
    if (!all) {
        hidehidden(cb_dirs);
    }

    // show file's inode
    if (inodes) {
        printdetails(cb_dirs);
    } else {
        printnames(cb_dirs);
    }

    printf("\n");

    return 0;
}
