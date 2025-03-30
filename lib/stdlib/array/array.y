// std/array/array.lang

type Array {
    let internalData;

    // map(fn) -> new array
    function map(fn) -> {
        asm("builtin_map");
    }

    // forEach(fn)
    function forEach(fn) -> {
        asm("builtin_forEach");
    }
}