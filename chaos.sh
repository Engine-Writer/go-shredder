#!/bin/bash
set -euo pipefail

TARGET_DIR="./rand"
MAX_DEPTH=5
MAX_FILES_PER_DIR=8
MAX_SUBDIRS_PER_DIR=4
MAX_FILE_SIZE=8192

mkdir -p "$TARGET_DIR"

random_name() {
    tr -dc a-z0-9 </dev/urandom | head -c $((4 + RANDOM % 9))
}

create_random_tree() {
    local dir="$1"
    local depth="$2"

    # Random files
    local num_files=$((1 + RANDOM % MAX_FILES_PER_DIR))
    for ((i=0; i<num_files; i++)); do
        local file_path="$dir/$(random_name).bin"
        head -c $((1 + RANDOM % MAX_FILE_SIZE)) </dev/urandom >"$file_path"
        echo "Created file: $file_path"
    done

    # Random subdirs (force at least one if depth < MAX_DEPTH)
    if (( depth < MAX_DEPTH )); then
        local num_subdirs=$((1 + RANDOM % MAX_SUBDIRS_PER_DIR))
        for ((i=0; i<num_subdirs; i++)); do
            local subdir="$dir/$(random_name)"
            mkdir -p "$subdir"
            echo "Created dir: $subdir"
            create_random_tree "$subdir" $((depth + 1))
        done
    fi
}

create_random_tree "$TARGET_DIR" 1

echo "Chaotic random filesystem generated in $TARGET_DIR"
