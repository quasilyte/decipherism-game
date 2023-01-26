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

To build a game for wasm (browser):

```bash
GOOS=js GOARCH=wasm go build -o ../web/main.wasm .
```

After that, a `web` folder will contain 3 files:

* index.html
* wasm_exec.js
* main.wasm

Put these files into a single archive to create an itch-io uploadable bundle.

This game is tested on these targets:

* windows/amd64
* linux/amd64
* js/wasm
