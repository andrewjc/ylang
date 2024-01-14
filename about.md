# YLang Language Specification

## Introduction
YLang merges static typing and object-oriented paradigms with features of functional programming. Its design aims for succinctness, expressiveness, and a developer-centric approach, drawing inspiration from various contemporary languages.

## Syntax Guidelines

### Naming Rules
- **Methods & Variables**: Adhere to camelCase.
- **Classes & Namespaces**: Follow PascalCase.

### Data Types
- **Static Typing**: Infer types statically, with optional explicit declarations.
- **Primitive Types**: Includes `int`, `float`, `bool`, `string`, `char`.
- **Collections**: Implements `List<T>`, `Set<T>`, `Map<K, V>`.
- **Pointers**: Utilizes pointer syntax `Type*`.

### Defining Classes
```plaintext
class ClassName {
    // Attributes and Methods
}
```

### Methods
- **Static Methods**: Defined as `static returnType methodName(params) -> body`.
- **Instance Methods**: Follow `returnType methodName(params) -> body`.
- **Lambdas**: Use `(params) -> expression` or extended block `(params) -> { body }`.

### Constructors
- Constructors implicitly return the object instance.
- Defined using `constructor(params) -> body`.

### Object Creation and Method Chaining
- Objects are created via `let objectName(ClassName, params)`.
- Enable method chaining with `->`.

### Control Flow Constructs
- **Conditional**: Standard if-else constructs.
- **Iteration**: Includes `for`, `while`, and collection-based `for item in collection`.
- **Switch-Case**: Utilize pattern matching with `switch`.

### Exception Handling
- Implements `try`, `catch`, `finally` for error management.

### Functional Features
- Implements `map`, `filter`, `reduce` for collections.
- Supports lambdas and higher-order functions.

### Lifecycle Hooks
- **onConstruct**: Defined as `onConstruct(lambdaAction) -> { lambdaAction(this); return this; }`.
- **onDestruct**: Specified with `onDestruct(lambdaAction) -> { /* Registration logic */ }`.

## Standard Library
- Includes basic IO, networking, and file operations.
- Provides standard data structures and algorithms.

## Language Integration
- Offers interoperability mechanisms with languages like C, Java.

## Memory Management
- Features automatic garbage collection.
- Allows manual memory management for expert users.

## Example
```plaintext
class Sample {
    let value;

    constructor(value) -> {
        this.value = value;
    }

    method doSomething() -> {
        print("Value is: " + this.value);
    }
}

main() -> {
    let sample(Sample, 42) -> onConstruct(lambda (x) -> print("Sample created"));
    sample.doSomething();
}
```

### Language Syntax Specification:

```
identifier ::= letter (letter | digit)*
number ::= digit+ ('.' digit+)?
string ::= '"' character* '"'
character ::= <any Unicode character except '"'>
letter ::= [a-zA-Z_]
digit ::= [0-9]

expression ::= term (('+' | '-') term)* | ternaryExpression
term ::= factor (('*' | '/') factor)*
factor ::= number | identifier | '(' expression ')'

ternaryExpression ::= traditionalTernary | arrowStyleTernary | colonPrefixedTernary | lambdaStyleTernary | inlineIfElseTernary
traditionalTernary ::= expression '?' expression ':' expression
arrowStyleTernary ::= expression '->' expression

':' expression
colonPrefixedTernary ::= expression ':' expression '?' expression
lambdaStyleTernary ::= '(' expression ')' '->' '{' expression '}' ':' '{' expression '}'
inlineIfElseTernary ::= 'if' expression 'then' expression 'else' expression

statement ::= variableDeclaration | functionCall | assignment | controlStatement | assemblyStatement
variableDeclaration ::= 'let' identifier ('(' typeName ')' )? '=' expression
functionCall ::= identifier '(' argumentList? ')'
assignment ::= identifier '=' expression
controlStatement ::= ifStatement | forStatement | whileStatement | doStatement | switchStatement

function ::= 'function' identifier '(' parameterList? ')' (':' returnType)? '->' block
returnType ::= typeName
lambda ::= '(' parameterList? ')' '->' (expression | block)
parameterList ::= parameter (',' parameter)*
parameter ::= identifier (':' typeName)?
argumentList ::= expression (',' expression)*
block ::= '{' statement* '}'

classDeclaration ::= classLambdaStyle | classTypeStyle
classLambdaStyle ::= identifier '=>' '{' classMember* '}'
classTypeStyle ::= 'type' identifier '{' classMember* '}'

dataStructure ::= dataBraces | dataEquals | dataColon | tupleLike
dataBraces ::= 'data' identifier '{' fieldList '}'
dataEquals ::= identifier '=' 'data' '{' fieldList '}'
dataColon ::= 'data' identifier ':' '{' fieldList '}'
tupleLike ::= identifier '=' '{' fieldList '}'

fieldList ::= 'let' field (',' field)*
field ::= identifier
classMember ::= variableDeclaration | methodDeclaration
methodDeclaration ::= (returnType? identifier '(' parameterList? ')' '->' block)

ifStatement ::= ifClassic | ifLambda
ifClassic ::= 'if' '(' expression ')' block ('else' block)?
ifLambda ::= 'if' lambda block ('else' lambda block)?

forStatement ::= forClassic | forLambda | forEach | forEachLambda
forClassic ::= 'for' '(' (variableDeclaration | assignment)? ';' expression ';' assignment ')' block
forLambda ::= 'for' lambda block
forEach ::= 'for' identifier 'in' 'range' '(' expression ',' expression ')' block
forEachLambda ::= 'for' identifier 'in' 'range' '(' expression ',' expression ')' '->' lambda

whileStatement ::= whileClassic | whileLambda
whileClassic ::= 'while' '(' expression ')' block
whileLambda ::= 'while' lambda block

doStatement ::= doClassic | doLambda
doClassic ::= 'do' block 'while' '(' expression ')'
doLambda ::= 'do' block 'while' lambda

switchStatement ::= 'switch' '(' expression ')' '{' switchCaseOrExpression* undefinedOrDefaultCase '}'
switchCaseOrExpression ::= ('case'? expression ':' block)+
undefinedOrDefaultCase ::= ('undefined' lambda | 'default') ':' block

onConstruct ::= 'onConstruct' lambda
onDestruct ::= 'onDestruct' lambda

program ::= mainFunction (classDeclaration | function | dataStructure)*
mainFunction ::= 'main' '()' '->' block

assemblyStatement ::= 'asm' block

assemblyBlock ::= 'asm' '{' assemblyCode '}'
assemblyCode ::= stringLiteral

comment ::= singleLineComment | multiLineComment
singleLineComment ::= '//' character* <end-of-line>
multiLineComment ::= '/*' character* '*/'
```