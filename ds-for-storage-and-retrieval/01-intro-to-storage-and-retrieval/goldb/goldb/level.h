#include "leveldb/c.h"

leveldb_t *initlevel();
void closelevel(leveldb_t *db);
int putlevel(leveldb_t *db, char *key, char *value);
char *getlevel(leveldb_t *db, char *key);