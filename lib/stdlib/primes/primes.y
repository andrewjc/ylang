// stdlib/primes - prime number utilities
// Implemented entirely in Y-lang using Linux x86-64 syscalls.
// No external C runtime is required.
//
// Syscall numbers (x86-64 Linux):
//   SYS_write  = 1
//   SYS_mmap   = 9
//   SYS_munmap = 11
//
// mmap flags:
//   PROT_READ|PROT_WRITE = 3
//   MAP_PRIVATE|MAP_ANONYMOUS = 34

// isPrime returns 1 if n is a prime number, 0 otherwise.
function isPrime(n) -> {
    if (n < 2) {
        return 0;
    }
    let i = 2;
    while (i * i <= n) {
        let q = n / i;
        let rem = n - q * i;
        if (rem == 0) {
            return 0;
        }
        i = i + 1;
    }
    return 1;
}

// print_digit writes a single ASCII digit d (0..9) to stdout (fd 1).
function print_digit(d) -> {
    let buf = syscall(9, 0, 8, 3, 34, -1, 0);
    buf[0] = d + 48;
    syscall(1, 1, buf, 1);
    syscall(11, buf, 8, 0, 0, 0, 0);
    return 0;
}

// print_int writes the decimal representation of non-negative integer n to stdout.
function print_int(n) -> {
    if (n < 10) {
        print_digit(n);
        return 0;
    }
    print_int(n / 10);
    let last = n - (n / 10) * 10;
    print_digit(last);
    return 0;
}

// print_primes prints the first count prime numbers to stdout, one per line.
function print_primes(count) -> {
    let found = 0;
    let n = 2;
    while (found < count) {
        if (isPrime(n) == 1) {
            print_int(n);
            syscall(1, 1, "\n", 1);
            found = found + 1;
        }
        n = n + 1;
    }
    return 0;
}
