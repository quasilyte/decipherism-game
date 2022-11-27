# Building from Source

This should be enough to build the game from sources:

```bash
cd src
go build -o ../bin/decipherism .
```

> You will need a go 1.18+ in order to build this game.

If you want to build a game for a different platform, use Go cross-compilation:

```bash
GOOS=windows go build -o ../bin/decipherism.exe .
```

This game is tested on these targets:

* windows/amd64
* linux/amd64
* js/wasm
