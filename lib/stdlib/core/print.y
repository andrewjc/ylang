// std/core/print.y

function print(x) -> {
    let len = 0;
    let p = str;
    while (p* != 0) {
        p = p + 1;
        len = len + 1;
    }
    syscall(1, 1, str, len, 0, 0, 0);
}
