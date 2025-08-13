# Go Shredder

This project provides a **fast and safe file shredding utility** written in Go. It overwrites files with cryptographically secure random data multiple times and then deletes them, handling files, directories, and symlinks correctly.

This utility can recursively shred directories, skip errors (given `force=true`), and ensures symlinks are deleted safely without following them to avoid accidental corruption of critical symbolically-linked files.

---

## Project Overview

* **Shred Function**: Overwrites a file multiple times with random data and deletes it.
* **Recursive Deletion**: Directories can be shredded recursively, with optional concurrency limits.
* **Symlink Safety**: Symlinks are deleted without touching their target.
* **Concurrency**: Uses goroutines and a semaphore to speed up shredding of large directories safely.
* **Command-line Interface**: Configurable iterations, recursion, and error handling via flags.

---

## Key Features

* Shred normal files -> overwrite + delete
* Shred directories -> recursively overwrite all contents + delete
* Shred symlinks -> delete without touching the target
* Configurable overwrite iterations
* Optional recursive mode
* Optional "force" mode to continue on errors
* Asyncronous file handling to speed up shredding large directories
* Chunked writes for large files to save memory

---

## Directory Structure

```
/go-shredder
|
+-- main.go                 # Main Go shredder implementation
+-- test.sh                 # Bash script for testing shredding functionality
+-- go.mod/                 # Go Modules required for project
+-- test_shred/             # Test environment folder (created by build.sh)
+-- README.md               # Project documentation (This file)
```

---

## How to Use the Shredder

Run the Go program with flags:

```
go run main.go -path="./target_file_or_dir" -iters=3 -recurse=true -force=true
```

### Command-line Flags:

| Flag       | Description                                    |
| ---------- | ---------------------------------------------- |
| `-path`    | File or directory to shred (required)          |
| `-iters`   | Number of overwrite iterations (default: 3)    |
| `-recurse` | Recursively shred directories (default: false) |
| `-force`   | Continue on errors (default: false)            |

Example: Shred a folder recursively with 3 overwrite passes and force errors:

```
go run main.go -path="./test_shred" -iters=3 -recurse=true -force=true
```

---

## Testing

A `test.sh` script sets up a test environment with files, directories, and symlinks, then runs the shredder. After completion, it confirms that all targets are removed.

```
sh test.sh
```

---

## Dependencies

* Go 1.20+ (or latest stable release)

Optional for testing:

* Bash shell
* `dd` for generating test files in `test.sh`

---

## Possible Use Cases

* Securely delete sensitive files from disk.
* Remove temporary files containing confidential data.
* Clear old backups safely before disposal.

### Advantages

* Handles many edge cases: files, directories, and symlinks.
* Fast concurrent processing for large directories.
* Simple CLI with flexible options.

### Drawbacks

* Does not handle read-only or locked files.
* Cannot guarantee secure deletion on networked or non-standard file systems.
* Might start too many Gorotines if you have a long directory chain (e.g. ./path1/path2/path3/path4/..../path999/) despite such chains are not possible in most consumer filesystems

---

## License

This project is licensed under no license because ~~I have no clue which license to pick~~ this is an interview project. The code may be used freely anywhere aslong as the file is equiped with a visible disclaimer linking to this GitHub repository