// Represents an array type. The actual data structure and memory management
// are handled natively by the compiler's code generation for array literals
// and operations.
type Array {

    let internalData;

    // function map(callback) -> Array
    // Applies the callback function to each element of the array
    // and returns a *new* array containing the results.
    // The implementation is provided natively by the compiler.
    function map(callback) -> {
        // Native implementation hook - body is ignored by the compiler
        // for recognized built-in/native methods like 'map'.
        // The return type (implicitly Array) should be handled by the native logic.
    }

    // function forEach(callback) -> void
    // Applies the callback function to each element of the array.
    // Does not return a new array (typically returns void or the original array).
    // The implementation is provided natively by the compiler.
    function forEach(callback) -> {
        // Native implementation hook - body is ignored by the compiler
        // for recognized built-in/native methods like 'forEach'.
    }

    // Add other potential array methods here (e.g., length, push, pop)
    // function length() -> int { /* native */ }
}