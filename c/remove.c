#include <stdio.h>
#include <string.h>

void removehidden(char **);

void removehidden(char **arr) {
    char* starts = ".";
    printf("%10s\n", arr[0]);
    printf("%10s\n", arr[1]);
    /* while (arr != NULL) { */
        /* if (strncmp(arr, starts, 0) == 0) { */
            /* printf("oh snap\n"); */
        /* } */
        /* arr++; */
    /* } */
}

int main() {
    char* names[3] = {"Abe", "Bonnie", ".Clyde"};

    removehidden(names);
}
