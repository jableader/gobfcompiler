#bfukt
####*A silly little language that compiles into a sillier, little one*

Have you ever noticed that most the languages that compile to brainfuck dont
have the same neat little heart that brainfuck does? They're either far too bf
still, or they are languages like C where the BF output looks nothing like
lovingly crafted brainfuck.

I have tried to design a language and compiler (this is V2 actually!) that can
create neat little bf code that looks like it could have been handwritten, but
is still at a higher level than simple macros.

##To Build the Compiler
Simply run `go build` in the project root. You may also use `go test` to run the
available tests.

##Usage
Simply `go build` the project. You have a few command line options such as `-lex`
to only print the lexicons and `-str` to print an 'assembly-like' view.

##The Language
The language is dynamically typed, although that said, there is only two types.
The two types are functions and variables, and they can be used interchangeably
(first class functions).

The language, as of V2 no longer has significant whitespace. I decided to take
this approach because I felt that the whitespacing added unnecessary complexity
to both the language and its compiler.

The spec below includes features that are not yet implemented. I will note when
a feature is not yet ready. You'll find other examples in the /examplecode
directory.

###Comments
A line that starts with a `#hash` is a comment

###Declaring a variable
Declaring a variable is simple, and all of the following are valid and do what
you would expect them to. Names are prefixed with a `$`.

```
var $j;
var $foo, $bar;

#Combined declaration and assignment is yet to be implemented

# Declare foo, set to 5
var $foo = 5;

# Declare a & b, they will both have the ASCII value of 'c'
var a, b = 'c'

# Declare c, and set it's value to be the same as a
var c = a
```

### Arithmatic
This is where things get a little interesting. You may set a variable to have
the value of a literal, a character or another variable. If the variable is
prefixed by an underscore, then that tells the compiler that its value may
be zero after the operation. This allows much smaller brainfuck code.

You may add, subtract and set values.

```
# Set a to 5
$a = 5;

# Add 5 to a
+$a = 5;

# Subtract 5 from a
-$a = 5;

# Set a to 'A'
$a = 'A';

# Set a to b, with b keeping its original value
$a = $b;

# Set a to b, leaving b at zero
$a = _$b;

# Add 5 to a and b
+$a, +$b = 5;

# Add c to a, subtract c from b and leave c at zero
+$a, -$b = _$c;
```

### IO
IO is pretty simple in the language. Simply use the command `print` to print,
the operation `read` to read input.

```
# Read input into a
read $a;

# Read input into a and then into b (two reads)
read $a, $b;

# Print a
print $a

# Print a and b
print $a, $b
```

### If Statement
The language offers the all important `if` statement, which means 'if this
variable is greater than zero'. You also have an else case. Just like arithmetic
it produces much smaller code if the variable is allowed to be zero afterwards,
so where possible remember to put in that underscore.

Note, you should not refer to the conditional parameter inside your method,
crazy things may happen, like infinite loops.

```
# If a is greater than zero, add 1 to b (leave a at zero)
if _$a {
	+$b = 1;
}

# If a is greater than zero, add 1 to b (leave a at its original value)
if $a {
	+$b = 1;
}

# If a is greater than zero, add 1 to b, otherwise subtract 2 from b
if _$a {
	+$b = 1;
} else {
	-$b = 2
}
```

### Loops!
The language offers only the simple while loop, and it simply means 'while the
variable is not zero'. You can also precede your variable with a `-` so that it
automatically decrements at the end of the loop.

```
# Simple loop
while $a {
	print $b;
	-$a = 1;
}

# Functionally identical loop
while -$a {
	print $b;
}
```

### Functions
The language offers simple functions. Your variables are passed by reference and
the functions do not have return values. They effectively operate similarly to
macros.

You may also pass functions around as parameters to other functions.

```
# Simple add function
def $add($v) {
	+$v = 5;
}

var $a = 0;
$add($a);

# Function to call another function five times
def doFiveTimes($f, $v) {
	var $n = 5;
	while -$n {
		$f($v);
  }

var $b = 0;
$doFiveTimes($add, $b);
```
