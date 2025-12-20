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



Day 5:
Implemented poor man's objects + field access. Setup tiny testing framework that executes some sample programs and matches its print statements output to expected output stored in comments (`#? STDOUT expected_value`). This doesn't quite test cases that are supposed to fail though. Something to look into later.

Implement boolean not operator. Begin hooking up Raylib functions (each function is its own instruction in the bytecode). We can open windows now woo.

Implemented ability to call functions without necessarily assigning them (e.g `_ = a()` => `a()`). Implement constant globals as well as a way to read them. So far only INT/Float globals are supported and are converted directly into IntValue/FloatValue objects accordingly.

Day 6:
Parser-level hack to allow negative constants. 

Implement time() and rlIsKeyPressed() operations. Begin actually working on the game.

Implemented casting to float / greater equal operator.

Implemented player grid-based movement.


Implemented World class, which essentially is a collection of columns containing crates.

Implemented collision checking while moving the player horizontally. Given player's new_goal_x, we sample the column at index `new_goal_x` and check if any of the crates overlap with a player along the Y-axis. This overlap check is handled by `Interval` object and `intervals_overlap` function. 


Implemented player jumping by introducing `goal_y` field. When the user requests jump, this `goal_y` is set to be current `y` position + 1, and `PLAYER_JUMPING` state is set.

Gravity and jumping is simulated in `player_tick_vertical_movement`  function. Whoops, we are not taking column's crate position into account while applying gravity.

World object now has method `world_get_highest_in_column_below` which samples the highest ground level that is below provided y coordinate. Implemented naively without taking a scenario of player sinking into the crate into account. 

So now, when gravity is applied, sample player's current ground level using `world_get_highest_in_column_below` and if player is moving then we sample the target cell as well and take the maximum of these two ground heights. Gravity seems to work as expected.


Implemented random_int/random_float built-ins. We can now generate random colors for crates and we can now spawn crates`random_float() < 0.0008`. Also, work on applying crate gravity. I realized that object comparison doesn't quite work, should work based on reference comparison. Something to look into after the deadline. For crate gravity, we find the crate's column, take top-y coordinate of the crate underneath it (if present) and move the crate up until that point.


Day 7: 
Pushing crates. Correct gravity application for falling crates.


Bugfix: had evaluation stack overflow happening after running the game for a few minutes. Few `printf` statements in the interpreter later, I found out that stack top increased while spawning crates. Commenting out lines in that function, `self.columns[x] = column + crate;` was the culprit. Surely enough, `subscriptSet` had a bug: I was pushing a popped value onto the stack again after setting value in the array. And that was the fix. Quite surprized that the game has behaved in expected way.
