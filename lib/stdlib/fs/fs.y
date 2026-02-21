// stdlib/fs - filesystem operations

// listdir lists all files in the current working directory, printing each name to stdout.
function listdir() -> {
    asm("builtin_list_cwd");
}
