// std/core/print.y

// print writes a string value to stdout.
function print(str) -> {
    asm("builtin_print_str", str);
}
