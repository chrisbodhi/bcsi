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