// stdlib/fs/storage_service.y
// Storage Service implementation
//
// Filesystem and persistent storage service running in user space.
// Implemented entirely in Y-lang via Linux x86-64 syscalls.
// No external C runtime is required.
//
// Syscall numbers (x86-64 Linux):
//   SYS_read      =   0
//   SYS_write     =   1
//   SYS_close     =   3
//   SYS_lseek     =   8
//   SYS_mmap      =   9
//   SYS_munmap    =  11
//   SYS_rename    =  82
//   SYS_mkdir     =  83
//   SYS_unlink    =  87
//   SYS_getdents64= 217
//   SYS_openat    = 257
//
// Constants:
//   AT_FDCWD       = -100   (open relative to cwd)
//   O_RDONLY        =    0
//   O_WRONLY        =    1
//   O_RDWR          =    2
//   O_CREAT         =   64   (0x40)
//   O_TRUNC         =  512   (0x200)
//   O_APPEND        = 1024   (0x400)
//   SEEK_SET        =    0
//   SEEK_CUR        =    1
//   SEEK_END        =    2
//   PROT_READ       =    1
//   PROT_WRITE      =    2
//   MAP_PRIVATE     =    2
//   MAP_ANONYMOUS   =   32   (0x20)
//   DEFAULT_MODE    =  420   (0644 octal)

import "stdlib/core";

// --------------------------------------------------------------------------
// Memory helpers
// --------------------------------------------------------------------------

// allocBuffer allocates a scratch buffer of the given size via mmap(2).
// Returns the buffer address, or a value < 0 on failure.
function allocBuffer(size) -> {
    // mmap(addr=0, length=size, prot=PROT_READ|PROT_WRITE=3,
    //      flags=MAP_PRIVATE|MAP_ANONYMOUS=34, fd=-1, offset=0)
    let buf = syscall(9, 0, size, 3, 34, -1, 0);
    return buf;
}

// freeBuffer releases a buffer previously obtained from allocBuffer.
function freeBuffer(buf, size) -> {
    // munmap(buf, size)
    syscall(11, buf, size, 0, 0, 0, 0);
}

// --------------------------------------------------------------------------
// Low-level file operations
// --------------------------------------------------------------------------

// storageOpen opens a file relative to the current working directory.
// flags: O_RDONLY=0, O_WRONLY=1, O_RDWR=2, combined with O_CREAT=64, etc.
// mode: permission bits (e.g. 420 for 0644).
// Returns the file descriptor, or a value < 0 on error.
function storageOpen(path, flags, mode) -> {
    // openat(AT_FDCWD=-100, path, flags, mode)
    let fd = syscall(257, -100, path, flags, mode, 0, 0);
    return fd;
}

// storageClose closes a previously opened file descriptor.
// Returns 0 on success, or a value < 0 on error.
function storageClose(fd) -> {
    // close(fd)
    let ret = syscall(3, fd, 0, 0, 0, 0, 0);
    return ret;
}

// storageRead reads up to count bytes from fd into buf.
// Returns the number of bytes actually read, 0 at end-of-file,
// or a value < 0 on error.
function storageRead(fd, buf, count) -> {
    // read(fd, buf, count)
    let n = syscall(0, fd, buf, count, 0, 0, 0);
    return n;
}

// storageWrite writes count bytes from buf to fd.
// Returns the number of bytes written, or a value < 0 on error.
function storageWrite(fd, buf, count) -> {
    // write(fd, buf, count)
    let n = syscall(1, fd, buf, count, 0, 0, 0);
    return n;
}

// storageSeek repositions the file offset of fd.
// whence: SEEK_SET=0 (absolute), SEEK_CUR=1 (relative), SEEK_END=2.
// Returns the resulting offset from the beginning of the file,
// or a value < 0 on error.
function storageSeek(fd, offset, whence) -> {
    // lseek(fd, offset, whence)
    let pos = syscall(8, fd, offset, whence, 0, 0, 0);
    return pos;
}

// --------------------------------------------------------------------------
// File-level convenience helpers
// --------------------------------------------------------------------------

// storageFileSize returns the size in bytes of the file referred to by fd,
// determined by seeking to the end and back.  Returns < 0 on error.
function storageFileSize(fd) -> {
    // Remember current position
    let cur = syscall(8, fd, 0, 1, 0, 0, 0);
    if (cur < 0) {
        return cur;
    }
    // Seek to end
    let end = syscall(8, fd, 0, 2, 0, 0, 0);
    if (end < 0) {
        return end;
    }
    // Restore original position
    syscall(8, fd, cur, 0, 0, 0, 0);
    return end;
}

// --------------------------------------------------------------------------
// Path helpers
// --------------------------------------------------------------------------

// buildPath concatenates basePath + "/" + key into a 4096-byte buffer.
// Returns the buffer address (caller must freeBuffer when done),
// or a value < 0 on allocation failure.
function buildPath(basePath, key) -> {
    let pathBuf = allocBuffer(4096);
    if (pathBuf < 0) {
        return pathBuf;
    }

    // Copy basePath
    let baseLen = 0;
    while (basePath[baseLen]) {
        pathBuf[baseLen] = basePath[baseLen];
        baseLen = baseLen + 1;
    }

    // Append '/'
    pathBuf[baseLen] = 47;
    let pos = baseLen + 1;

    // Append key
    let keyLen = 0;
    while (key[keyLen]) {
        pathBuf[pos] = key[keyLen];
        pos = pos + 1;
        keyLen = keyLen + 1;
    }

    // Null-terminate
    pathBuf[pos] = 0;

    return pathBuf;
}

// --------------------------------------------------------------------------
// Persistent key-value store
// --------------------------------------------------------------------------
// A minimal flat-file storage service.  Each key is stored as a separate
// file under a configurable base directory.  Values are raw byte sequences.

// storageInit creates the base directory for the storage service.
// Returns 0 on success, or a value < 0 on error.
function storageInit(basePath) -> {
    // mkdir(basePath, 0755 = 493)
    let ret = syscall(83, basePath, 493, 0, 0, 0, 0);
    // EEXIST (-17) is acceptable — the directory may already exist.
    if (ret < 0) {
        if (ret == -17) {
            return 0;
        }
        return ret;
    }
    return 0;
}

// storagePut writes value (valueLen bytes) under the given key.
// The key is used directly as the filename inside basePath.
// Returns the number of bytes written, or < 0 on error.
function storagePut(basePath, key, value, valueLen) -> {
    let pathBuf = buildPath(basePath, key);
    if (pathBuf < 0) {
        return -1;
    }

    // Open (create/truncate) the file
    // O_WRONLY | O_CREAT | O_TRUNC = 1 | 64 | 512 = 577
    let fd = syscall(257, -100, pathBuf, 577, 420, 0, 0);
    if (fd < 0) {
        freeBuffer(pathBuf, 4096);
        return fd;
    }

    // Write value
    let written = syscall(1, fd, value, valueLen, 0, 0, 0);

    // Close
    syscall(3, fd, 0, 0, 0, 0, 0);

    freeBuffer(pathBuf, 4096);
    return written;
}

// storageGet reads the value stored under key into the caller-provided buf
// (up to bufLen bytes).
// Returns the number of bytes read, or < 0 on error (e.g. -2 for ENOENT
// when the key does not exist).
function storageGet(basePath, key, buf, bufLen) -> {
    let pathBuf = buildPath(basePath, key);
    if (pathBuf < 0) {
        return -1;
    }

    // Open for reading
    let fd = syscall(257, -100, pathBuf, 0, 0, 0, 0);
    if (fd < 0) {
        freeBuffer(pathBuf, 4096);
        // Return the syscall error (e.g. -2 for ENOENT) so callers can
        // distinguish "not found" from other failures.
        return fd;
    }

    // Read up to bufLen bytes
    let n = syscall(0, fd, buf, bufLen, 0, 0, 0);

    syscall(3, fd, 0, 0, 0, 0, 0);
    freeBuffer(pathBuf, 4096);

    if (n < 0) {
        return n;
    }
    return n;
}

// storageDelete removes the file corresponding to key from the store.
// Returns 0 on success, or < 0 on error.
function storageDelete(basePath, key) -> {
    let pathBuf = buildPath(basePath, key);
    if (pathBuf < 0) {
        return -1;
    }

    // unlink(pathBuf)
    let ret = syscall(87, pathBuf, 0, 0, 0, 0, 0);

    freeBuffer(pathBuf, 4096);
    return ret;
}

// storageList prints every key (file name) in the store to stdout.
// Omits the "." and ".." pseudo-entries.
// Returns 0 on success, or < 0 on error.
function storageList(basePath) -> {
    let buf = allocBuffer(4096);
    if (buf < 0) {
        return -1;
    }

    // Open the base directory
    let fd = syscall(257, -100, basePath, 0, 0, 0, 0);
    if (fd < 0) {
        freeBuffer(buf, 4096);
        return fd;
    }

    // getdents64(fd, buf, 4096)
    let nbytes = syscall(217, fd, buf, 4096, 0, 0, 0);

    let pos = 0;
    while (pos < nbytes) {
        // d_reclen at offset 16 (low byte sufficient for small dirs)
        let reclen = buf[pos + 16];

        // d_name starts at offset 19
        let nameAddr = buf + pos + 19;

        // Measure name length
        let nameLen = 0;
        while (nameAddr[nameLen]) {
            nameLen = nameLen + 1;
        }

        // Skip "." and ".."
        let isDot = 0;
        if (nameLen == 1) {
            if (nameAddr[0] == 46) {
                isDot = 1;
            }
        }
        if (nameLen == 2) {
            if (nameAddr[0] == 46) {
                if (nameAddr[1] == 46) {
                    isDot = 1;
                }
            }
        }

        if (isDot == 0) {
            // Print key name followed by newline
            syscall(1, 1, nameAddr, nameLen, 0, 0, 0);
            syscall(1, 1, "\n", 1, 0, 0, 0);
        }

        pos = pos + reclen;
    }

    syscall(3, fd, 0, 0, 0, 0, 0);
    freeBuffer(buf, 4096);
    return 0;
}

// storageRename renames a key in the store (moves the underlying file).
// Returns 0 on success, or < 0 on error.
function storageRename(basePath, oldKey, newKey) -> {
    let oldPath = buildPath(basePath, oldKey);
    if (oldPath < 0) {
        return -1;
    }
    let newPath = buildPath(basePath, newKey);
    if (newPath < 0) {
        freeBuffer(oldPath, 4096);
        return -1;
    }

    // rename(oldPath, newPath)
    let ret = syscall(82, oldPath, newPath, 0, 0, 0, 0);

    freeBuffer(oldPath, 4096);
    freeBuffer(newPath, 4096);
    return ret;
}

// storageAppend appends valueLen bytes of value to the file for key.
// Creates the file if it does not exist.
// Returns the number of bytes written, or < 0 on error.
function storageAppend(basePath, key, value, valueLen) -> {
    let pathBuf = buildPath(basePath, key);
    if (pathBuf < 0) {
        return -1;
    }

    // O_WRONLY | O_CREAT | O_APPEND = 1 | 64 | 1024 = 1089
    let fd = syscall(257, -100, pathBuf, 1089, 420, 0, 0);
    if (fd < 0) {
        freeBuffer(pathBuf, 4096);
        return fd;
    }

    let written = syscall(1, fd, value, valueLen, 0, 0, 0);

    syscall(3, fd, 0, 0, 0, 0, 0);
    freeBuffer(pathBuf, 4096);
    return written;
}
