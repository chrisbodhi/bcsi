/*
 * Minimal rssfs example for macFUSE 4.0.5, adapted from
 * https://github.com/libfuse/libfuse/blob/master/example/hello.c
 *
 * Distributed under GNU GPLv2 per the original
 *
 * Compile:
 *
 *     clang -Wall rssfs.c `pkg-config fuse --cflags --libs` -o rssfs
 *
 * Run:
 *     mkdir mnt-point
 *     ./rssfs mnt-point/
 *
 * ... then cat mnt-point/hello to see "Hello World!"
 *
 */

#define FUSE_USE_VERSION 26

#include <fuse.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <stddef.h>
#include <assert.h>

/*
 * Command line options
 *
 * We can't set default values for the char* fields here because
 * fuse_opt_parse would attempt to free() them when the user specifies
 * different values on the command line.
 */
static struct options {
	const char *filename;
	const char *contents;
	const char *feed;
	int show_help;
} options;

#define OPTION(t, p)                           \
    { t, offsetof(struct options, p), 1 }
static const struct fuse_opt option_spec[] = {
	OPTION("--name=%s", filename),
	OPTION("--contents=%s", contents),
	OPTION("--feed=%s", feed),
	OPTION("-h", show_help),
	OPTION("--help", show_help),
	FUSE_OPT_END
};

static int rssfs_getattr(const char *path, struct stat *stbuf) {
  int res = 0;

  // check stat(2) for these constants
  memset(stbuf, 0, sizeof(struct stat));
  // if we're at the root...
  if (strcmp(path, "/") == 0) {
	// ...then update root's directory's stat's struct...
	// for file stat's permission bits, and...
    stbuf->st_mode = S_IFDIR | 0755;
	// ...number of hard links
    stbuf->st_nlink = 2;
  // ...if we're at our one file, which is named hello, by default...
  } else if (strcmp(path + 1, options.filename) == 0) {
	// ...update the file stat's permission bits, and...
    stbuf->st_mode = S_IFREG | 0444;
	// ...number of hard links, and...
    stbuf->st_nlink = 1;
	// ...file size, based on contents
    stbuf->st_size = strlen(options.contents);
  } else
  	// ...else, update the res to return a non-zero exit code
  	// error: no file or directory
    res = -ENOENT; // this will be -2

  return res;
}

static int hello_readdir(const char *path, void *buf, fuse_fill_dir_t filler,
                         off_t offset, struct fuse_file_info *fi) {
  (void)offset;
  (void)fi;

  if (strcmp(path, "/") != 0)
    return -ENOENT;

  filler(buf, ".", NULL, 0);
  filler(buf, "..", NULL, 0);
  filler(buf, options.filename, NULL, 0);

  return 0;
}

static int hello_open(const char *path, struct fuse_file_info *fi)
{
	if (strcmp(path+1, options.filename) != 0)
		return -ENOENT;

	if ((fi->flags & O_ACCMODE) != O_RDONLY)
		return -EACCES;

	return 0;
}

static int hello_read(const char *path, char *buf, size_t size, off_t offset,
		      struct fuse_file_info *fi)
{
	size_t len;
	(void) fi;
	if(strcmp(path+1, options.filename) != 0)
		return -ENOENT;

	len = strlen(options.contents);
	if (offset < len) {
		if (offset + size > len)
			size = len - offset;
		memcpy(buf, options.contents + offset, size);
	} else
		size = 0;

	return size;
}

static const struct fuse_operations hello_oper = {
	.getattr	= rssfs_getattr,
	.readdir	= hello_readdir,
	.open		= hello_open,
	.read		= hello_read,
};

static void show_help(const char *progname)
{
	printf("usage: %s [options] <mountpoint>\n\n", progname);
	printf("File-system specific options:\n"
	       "    --name=<s>          Name of the \"hello\" file\n"
	       "                        (default: \"hello\")\n"
	       "    --contents=<s>      Contents \"hello\" file\n"
	       "                        (default \"Hello, World!\\n\")\n"
     	   "    --feed=<s>      	RSS feed to explore\n"
	       "                        (default \"https://omar.website/posts/index.xml\")\n"
	       "\n");
}

int main(int argc, char *argv[])
{
	int ret;
	struct fuse_args args = FUSE_ARGS_INIT(argc, argv);

	/* Set defaults -- we have to use strdup so that
	   fuse_opt_parse can free the defaults if other
	   values are specified */
	options.filename = strdup("hello");
	options.contents = strdup("Hello World!\n");
	options.feed = strdup("https://omar.website/posts/index.xml");

	/* Parse options */
	if (fuse_opt_parse(&args, &options, option_spec, NULL) == -1)
		return 1;

	/* When --help is specified, first print our own file-system
	   specific help text, then signal fuse_main to show
	   additional help (by adding `--help` to the options again)
	   without usage: line (by setting argv[0] to the empty
	   string) */
	if (options.show_help) {
		show_help(argv[0]);
		assert(fuse_opt_add_arg(&args, "--help") == 0);
		args.argv[0][0] = '\0';
	}

	ret = fuse_main(args.argc, args.argv, &hello_oper, NULL);
	fuse_opt_free_args(&args);
	return ret;
}
