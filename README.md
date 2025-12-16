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
