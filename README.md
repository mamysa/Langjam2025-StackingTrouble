Day 1

Basic lexing / parsing simple arithmetic expressions / print + assign statements. Emitting bytecode for these features.

Day 2
More parsing - function definition syntax. Implementing basic bytecode interpreter. Direct function calls.

Day 3
Virtual function calls. Need to keep track of locals in the function to make a correct decision whether to call Virtual/Direct. Function names are not shadowed properly.

So things like 
```
def myfunc() {
print(0);
}

def a(n) {
print(n);

}

def main() {

myfunc = &a;
myfunc(n);  <-- should be CallVirtual BUT because myfunc is defined, Call <Offset> is done instead.
}
```
dont actually work.

TODO: should really split "real" instructions from IR instructions (that use textual labels instead of actual offsets).
Implemented if statements
Implemented while loops
Implemented ternary expressions.
Got rid of IrInstruction business. Was very annoying to deal with.


Day 4: 
Complete boolean and/or operations with short-circuiting.
a or b -> if a then a else b 
a and b -> if a then b else a

Implement ListValue type, ability to initialize lists, append items to list, array subscript syntax (on both sides of the assignment expression). Implement Len operator. Implement FloatValue (with casting int to float when appropriate). Ready for proper testing now. (probably will use python script that simply calls the go compiler with appropriate program).


