# Stacking Trouble
Stacking Trouble is a game that (poorly) attempts to re-create Stack Attack game that came preloaded on old Siemens phones. The goal of the game to organize falling crates into a row by pushing them. Once the bottom-most row is filled, it is discarded and points are awarded. Crates can also be destroyed by jumping and hitting them with the character from underneath. 

This game is made for the [Langjam Gamejam](https://langjamgamejam.com/), the goal of which is to make some programming language first and then write a game in it. Stacking Trouble is written in Bla - a very Python-inspired dynamically-typed imperative interpreted programming language. The compiler and the interpreter are written in Go. See `stacking-trouble-game.bla` for game's source code and see `tests` directory for (incomplete) testing suite.

The language supports:
* Functions/Function references (i.e. virtual calls)
* Lists/Objects.
* Loop/If statements.
* Global constants.

[!screenshot](https://github.com/mamysa/Langjam2025-StackingTrouble/blob/main/stacking-trouble-screenshot.png)

# Building/Running

This project is using Go 1.23.5 (not tested with newer versions) and `raylib-go. Binaries for macOS and Windows are provided on the Releases page. The Windows binary includes `start` batch file for convenience.

The game can be run `bla-run stacking-trouble.bla`. Additionally, there are a few other command line arguments.

* `DEBUG=tokens` - display tokens. Not particularly useful for larger files.
* `DEBUG=ast` - display abstact syntax tree. Not particularly useful for larger files either.
* `DEBUG=bytecode` - displays results of compilation that will be fed into VM. 

## Building

If you want to build from source:

MacOs:

```
go build
```

Windows: 
See [raylib-go](https://github.com/gen2brain/raylib-go) installation instructions.

```
CGO_ENABLED=1 go build
./compiler.exe stacking-trouble.bla
```


# Known Bugs, etc

* Doesn't implement Stack Attack fully - there were different crate power-ups (such as extra-health, TNT crate, etc) that I ran out of time to implement.
* Missing game art. It is all just randomly coloured-boxes.
* Compiler codebase is a bit of a mess and needs cleanup. 

Bug reports are welcome!

