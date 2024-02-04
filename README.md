# Y-Lang

An experimental programming language that steals features from all my favourite languages.

## 1. Function Definitions

### Complex Lambda Functions

```
let complexLambda = (x, y) -> {
    let result = x * y;
    return result;
};

// Usage of lambda in a function
function calculate(input) -> {
    return complexLambda(input, input + 10);
}
```

### Lambda in Ternary Expressions

```
let lambdaTernary = (condition) -> {
    return condition ? (x -> x * 2) : (x -> x / 2);
};

// Usage example
let result = lambdaTernary(true)(5);
```

### Inline If-Else Ternary

```
let inlineTernary = if x > 10 then (x -> x * 2) else (x -> x - 5);
let result = inlineTernary(12); // Should return 24
```

## 2. Class Declarations

### Class with Lambda Style Methods

```
type MyClass {
    let value = 10;

    increase = (amount) -> {
        this.value += amount;
    };

    decrease = (amount) -> {
        this.value -= amount;
    };
}
```

### Data Structures

```
data MyData {
    let attributeOne, let attributeTwo
};

let myDataInstance = MyData { attributeOne = "Value1", attributeTwo = "Value2" };
```

## 3. Control Structures

### Lambda in For Loops

```
for i in range(0, 10) -> lambda (i) {
    print("Value: " + i);
}
```

### If Statement with Lambda

```
if (x > 5) -> {
    print("x is greater than 5");
} else -> {
    print("x is less or equal to 5");
}
```

## 4. Special Constructs

### Resource Management with Lambdas

```
class ResourceHandler {
    constructor(resource) -> {
        this.resource = resource;
    }

    onConstruct(lambdaAction) -> {
        lambdaAction(this);
        return this;
    }

    onDestruct(lambdaAction) -> {
        lambdaAction(this);
    }

    useResource() -> {
        print("Using resource: " + this.resource);
    }
}

// Usage example
let handler = ResourceHandler("Resource1")
    -> onConstruct((x) -> print("Constructed with resource: " + x.resource))
    -> onDestruct((x) -> print("Destructing, releasing resource: " + x.resource));

handler.useResource();
```

## 5. Program Structure

Main Function with Complex Lambdas

```
main() -> {
    let process = (input) -> {
      return input * 2;
    };
    
    let values = [1, 2, 3, 4, 5];
    values.map(process).forEach(print);
}
```

## 6. Advanced Ternary Expressions

Arrow Style Ternary

```
let arrowTernary = x > 5 -> "Greater" : "Less or Equal";
```

Lambda Style Ternary

```
let complexTernary = (x > 5) -> { x * 2 } : { x / 2 };
```

## 7. Assembly Integration

Inline Assembly Statement

```
function performLowLevelOperation() -> {
    asm {
        "mov eax, 1" // Example assembly instruction
    }
}
```

## 8. Extended Lambda Usage

Lambda in Variable Declaration

```
let myLambda = (x, y) -> x + y;
let result = myLambda(5, 10); // Should return 15
```

Lambda in Control Statement

```
if ((x, y) -> x == y)(5, 5) {
    print("x and y are equal");
}
```

## 9. Data Structures

Tuple-like Data Structure

```
MyTuple = { let first, let second };
let myTupleInstance = MyTuple { first = "Hello", second = "World" };
```