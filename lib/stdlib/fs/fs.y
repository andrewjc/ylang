// stdlib/fs - filesystem operations
// Implemented entirely in Y-lang via Linux x86-64 syscalls.
// No external C runtime is required.
//
// Syscall numbers (x86-64 Linux):
//   SYS_write     =   1
//   SYS_close     =   3
//   SYS_mmap      =   9
//   SYS_munmap    =  11
//   SYS_getdents64= 217
//   SYS_openat    = 257
//
// Constants:
//   AT_FDCWD   = -100   (open relative to cwd)
//   O_RDONLY   =    0
//   PROT_READ  =    1
//   PROT_WRITE =    2
//   MAP_PRIVATE =   2
//   MAP_ANONYMOUS= 32   (0x20)
//
// linux_dirent64 layout (each entry in getdents64 buffer):
//   offset  0: d_ino    (8 bytes)
//   offset  8: d_off    (8 bytes)
//   offset 16: d_reclen (2 bytes, little-endian; low byte is enough for dirs < 4 KiB)
//   offset 18: d_type   (1 byte)
//   offset 19: d_name   (null-terminated string)

// listdir prints the name of every entry in the current working directory,
// one per line, to stdout.  The synthetic entries "." and ".." are omitted.
function listdir() -> {
    // Allocate a 4096-byte scratch buffer via mmap(2).
    // mmap(addr=0, length=4096, prot=PROT_READ|PROT_WRITE=3,
    //      flags=MAP_PRIVATE|MAP_ANONYMOUS=34, fd=-1, offset=0)
    let buf = syscall(9, 0, 4096, 3, 34, -1, 0);

    // MAP_FAILED is returned as -1 (0xFFFFFFFFFFFFFFFF as i64).
    // Since signed comparisons treat that as negative, buf < 0 detects failure.
    if (buf < 0) {
        return 1;
    }

    // Open the current directory.
    // openat(AT_FDCWD=-100, ".", O_RDONLY=0)
    let fd = syscall(257, -100, ".", 0, 0, 0, 0);

    // A negative fd means the open failed; clean up and bail out.
    if (fd < 0) {
        syscall(11, buf, 4096, 0, 0, 0, 0);
        return 1;
    }

    // Read directory entries into buf.
    // getdents64(fd, buf, 4096)
    let nbytes = syscall(217, fd, buf, 4096, 0, 0, 0);

    // Walk each linux_dirent64 record in the buffer.
    let pos = 0;
    while (pos < nbytes) {
        // d_reclen is a little-endian u16 at offset 16 within the entry.
        // For directories small enough to fit in a 4 KiB buffer the low byte
        // is sufficient (reclen <= 255).
        let reclen = buf[pos + 16];

        // d_name begins at offset 19; compute its start address.
        let name_addr = buf + pos + 19;

        // Measure the name length (inline strlen).
        let name_len = 0;
        while (name_addr[name_len]) {
            name_len = name_len + 1;
        }

        // Skip the pseudo-entries "." (len==1) and ".." (len==2) whose first
        // character is '.'.
        let is_dot = 0;
        if (name_len == 1) {
            if (name_addr[0] == 46) {
                is_dot = 1;
            }
        }
        if (name_len == 2) {
            if (name_addr[0] == 46) {
                if (name_addr[1] == 46) {
                    is_dot = 1;
                }
            }
        }

        if (is_dot == 0) {
            syscall(1, 1, name_addr, name_len);
            syscall(1, 1, "\n", 1);
        }

        pos = pos + reclen;
    }

    // close(fd)
    syscall(3, fd, 0, 0, 0, 0, 0);

    // munmap(buf, 4096)
    syscall(11, buf, 4096, 0, 0, 0, 0);

    return 0;
}
