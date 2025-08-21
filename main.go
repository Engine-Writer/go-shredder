package main

import (
    "crypto/rand"
    "flag"
    "fmt"
    "os"
    "sync"
)

// HYPERPARAMETERS. ADJUST FOR SPEED IF NEEDED BUT BEWARE OF RESOURCE USAGE SPIKES
const chunkSize = 4096 * 1024   // 4MiB
const maxGoroutines = 16        // max parallel shredders

func randomBytes(size int) ([]byte, error) {
    buf := make([]byte, size)
    _, err := rand.Read(buf)
    return buf, err
}

func overwriteFile(path string) error {
    file, err := os.OpenFile(path, os.O_WRONLY, 0)
    if err != nil {
        return err
    }
    defer file.Close()

    info, err := file.Stat()
    if err != nil {
        return err
    }

    size := info.Size()
    buf := make([]byte, chunkSize)

    for offset := int64(0); offset < size; offset += chunkSize {
        toWrite := chunkSize
        if offset+int64(chunkSize) > size {
            toWrite = int(size - offset)
        }

        _, err := rand.Read(buf[:toWrite])
        if err != nil {
            return err
        }

        _, err = file.WriteAt(buf[:toWrite], offset)
        if err != nil {
            return err
        }
    }

    return nil
}

func shred_expanded(path string, iters int, recurse bool, force bool) error {
    info, err := os.Lstat(path)
    if err != nil {
        return err
    }

    // Symlink handling: delete the link itself, never follow (Be destructive but not too much)
    if info.Mode()&os.ModeSymlink != 0 {
        return os.Remove(path)
    }

    // Directory handling
    if info.IsDir() {
        if !recurse {
            return fmt.Errorf("path is a directory, set recurse=true to shred")
        }

        entries, err := os.ReadDir(path)
        if err != nil {
            return err
        }

        var wg sync.WaitGroup
        sem := make(chan struct{}, maxGoroutines)
        errChan := make(chan error, len(entries))

        for _, entry := range entries {
            childPath := path + "/" + entry.Name()
            wg.Add(1)

            sem <- struct{}{} // acquire slot
            go func(p string) {
                defer wg.Done()
                defer func() { <-sem }() // release slot

                if err := shred_expanded(p, iters, recurse, force); err != nil && !force {
                    errChan <- err
                }
            }(childPath)
        }

        wg.Wait()
        close(errChan)

        for e := range errChan {
            return e
        }

        return os.Remove(path)
    }

    // Normal file handling
    for i := 0; i < iters; i++ {
        if err := overwriteFile(path); err != nil {
            return err
        }
    }

    return os.Remove(path)
}

// To strictly comply with the requirement
func Shred(path string) error {
    return shred_expanded(path, 3, false, false); // no -r or -f specified in the requirements so I set it to false
}

func main() {
    path := flag.String("path", "", "File or directory to shred")
    iters := flag.Int("iters", 3, "Number of overwrite iterations")
    recurse := flag.Bool("recurse", false, "Recursively shred directories")
    force := flag.Bool("force", false, "Continue on errors")

    flag.Parse()

    if *path == "" {
        fmt.Println("Please provide a path using -path")
        os.Exit(1)
    }

    err := shred_expanded(*path, *iters, *recurse, *force)
    if err != nil {
        fmt.Println("Failed to shred:", err)
    } else {
        fmt.Println("Shredding completed successfully")
    }
}
