// std/core/print.y
// Implements print() entirely in Y-lang via Linux syscalls.
// No external C runtime is required.

// strlen returns the number of bytes before the null terminator.
function strlen(str) -> {
    let i = 0;
    while (str[i]) {
        i = i + 1;
    }
    return i;
}

// print writes str followed by a newline to stdout (fd 1) via SYS_write (1).
function print(str) -> {
    syscall(1, 1, str, strlen(str));
    syscall(1, 1, "\n", 1);
}
