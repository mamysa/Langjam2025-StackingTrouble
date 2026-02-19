# Stacking Trouble
Stacking Trouble is a game that (poorly) attempts to re-create Stack Attack game that came preloaded on old Siemens phones. The goal of the game to organize falling crates into a row by pushing them. Once the bottom-most row is filled, it is discarded and points are awarded. Crates can also be destroyed by jumping and hitting them with the character from underneath. 

This game is made for the [Langjam Gamejam](https://langjamgamejam.com/), the goal of which is to make some programming language first and then write a game in it. Stacking Trouble is written in Bla - a very Python-inspired dynamically-typed imperative interpreted programming language. The compiler and the interpreter are written in Go. See `stacking-trouble-game.bla` for game's source code and see `tests` directory for (incomplete) testing suite.

The language supports:
* Functions/Function references (i.e. virtual calls)
* Lists/Objects.
* Loop/If statements.
* Global constants.

Version 2 of the game adds the following features:

* Power-up crates - health crate, double speed crate and explosive crate. 
* Difficulty mode selector.
* Ability to pause the game without exiting.
* Proper graphics and animations.
* A menu displaying controls as well as in-game "how to play" tutorial.

Bug reports are welcome!


![screenshot](https://github.com/mamysa/Langjam2025-StackingTrouble/blob/main/stacking-trouble-screenshot.png?raw=true)

# Building/Running

This project is using Go 1.23.5 (not tested with newer versions) and `raylib-go`. Binaries for macOS and Windows are provided on the Releases page. The Windows version includes `run-stacking-trouble` batch file for convenience, and MacOs version includes `run-stacking-trouble.sh`.

The game can be run `bla-run stacking-trouble.bla`. Additionally, there are a few other command line arguments.

* `DEBUG=tokens` - display tokens. Not particularly useful for larger files.
* `DEBUG=ast` - display abstact syntax tree. Not particularly useful for larger files either.
* `DEBUG=bytecode` - displays results of compilation that will be fed into VM. 

## Building
If you want to build from source:

MacOs:

```
# BUILD 
go build -o bla-run

# RUN
bla-run stacking-trouble-game/stacking-trouble.bla
```

Windows: 
Building on Windows requires `gcc`, which can be acquired through Msys2, for example. See [raylib-go](https://github.com/gen2brain/raylib-go) installation instructions.

```
# BUILD
CGO_ENABLED=1 CGO_LDFLAGS="-static-libgcc -static -lpthread" go build -a -ldflags "-H=windowsgui" -o bla-run

# RUN
./bla-run.exe stacking-trouble-game/stacking-trouble.bla
```

# License

Code is licensed under MIT. Art assets are licensed under CC-BY 4.0.

`ThaleahFat.ttf` font is created by Rick Hoppmann and is licensed under CC-BY 4.0. It can be acquired on [Itch.io](https://tinyworlds.itch.io/free-pixel-font-thaleah).


