#!/bin/bash

set -e

TESTDIR="./test_shred"

echo "Setting up test environment..."
rm -rf "$TESTDIR"
mkdir -p "$TESTDIR/dir1/dir2"

# Create some files
echo "Hello World" > "$TESTDIR/file1.txt"
dd if=/dev/urandom of="$TESTDIR/file2.bin" bs=1K count=10

# Nested files
echo "Nested file" > "$TESTDIR/dir1/file3.txt"
dd if=/dev/urandom of="$TESTDIR/dir1/dir2/file4.bin" bs=1K count=5

# Symlinks
ln -s "$TESTDIR/file1.txt" "$TESTDIR/link_to_file"
ln -s "$TESTDIR/dir1" "$TESTDIR/link_to_dir"

echo "Before shredding:"
ls -lR "$TESTDIR"

echo "Running Go shredder..."
go run main.go -path="$TESTDIR" -iters=3 -recurse=true -force=true

echo "After shredding:"
ls -lR "$TESTDIR" || echo "(All files removed)"

echo "Test complete!"
