// Note: Assumes 'int' maps to a standard integer type (e.g., i32 or i64)
//       and '*int' is a pointer to that type.
// Pointers are 8 bits wide, integers are 32 bits wide.
// TODO: Assumes 'new' keyword handles allocation which isn't implemented yet.
// TODO: Assumes basic control flow (if, while) and operators which isn't implemented yet properly.
type Array {
    let length: int;  // Number of elements
    let data: *int;   // Pointer to the first element

    // This is an internal constructor, not meant to be called directly
    // as denoted by the __ prefix. It's used to create an Array object.
    // TODO: Will need to add compiler support for this.
    function __Array(len: int, data_ptr: *int) -> {
        self.length = len;
        self.data = data_ptr;
    }

    function map(self: Array, fn) -> Array {
        let len = self.length;

        // Allocate memory for the new array's data elements
        // Using 'new Type[size]' syntax for dynamic array allocation
        let result_data: *int = new int[len];

        // Allocate the new Array struct/object to hold the result
        // Using 'new Type' syntax for object allocation
        let result: Array = new Array;
        result.length = len;
        result.data = result_data;

        // Handle empty array case (optional optimization, depends on 'new int[0]' behavior)
        // if (len == 0) { return result; }

        let i = 0;
        while (i < len) {
            // Read element from original array's data
            let element = self.data[i]; // Requires index read support (GEP + Load)

            // Apply the function passed as argument
            let transformed = fn(element); // Standard function call

            // Write transformed element to the new array's data
            result.data[i] = transformed; // Requires index write support (GEP + Store)

            // Increment loop counter
            i = i + 1;
        }

        return result; // Return the newly created Array object
    }

    function forEach(self: Array, fn) -> Array {
         let len = self.length;
         let i = 0;
         while (i < len) {
             let element = self.data[i]; // Index read
             fn(element); // Call function, ignore result
             i = i + 1; // Increment
         }
         return self; // Return original array instance for chaining
    }
}