#include <stdio.h>
#include <string.h>
#include "leveldb/c.h"

leveldb_t *initlevel()
{
    char *dbname = "/tmp/leveldb";
    char *err = NULL;
    leveldb_options_t *options = leveldb_options_create();
    return leveldb_open(options, dbname, &err);
}

void closelevel(leveldb_t *db)
{
    leveldb_close(db);
}

int putlevel(leveldb_t *db, char *key, char *value)
{
    char *err = NULL;
    leveldb_writeoptions_t *opts = leveldb_writeoptions_create();
    leveldb_put(db, opts, key, strlen(key), value, strlen(value), &err);
    if (err)
    {
        printf("%s\n", err);
        return 1;
    }
    return 0;
}

char *getlevel(leveldb_t *db, char *key)
{
    char *val;
    size_t val_len;
    char *err = NULL;
    leveldb_readoptions_t *opts = leveldb_readoptions_create();
    val = leveldb_get(db, opts, key, strlen(key), &val_len, &err);
    return val;
}