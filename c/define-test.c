#include <stdio.h>

#define NKEYS (sizeof keytab / sizeof keytab[0])
// same as
// #define NKEYS (sizeof keytab / sizeof(struct key))

struct key {
    char* name;
    int count;
};

struct key keytab[] = {
{"A", 0},
{"B", 0},
{"F", 0}
};

int main() {
    printf("len is %lu\n", NKEYS);
    return 0;
}
