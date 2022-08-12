void tofive(int *n);

void tofive(int *n) {
    while (*n < 5) {
        (*n)++;
    }
}

int main() {
    int i;
    int *ip;

    i = 0;
    ip = &i;

    printf("i was %d...\n", i);
    tofive(ip);
    printf("i is now %d.\n", i);
}
