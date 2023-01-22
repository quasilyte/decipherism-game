package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/quasilyte/decipherism-game/leveldata"
	"github.com/quasilyte/ge/tiled"
)

func main() {
	log.SetFlags(0)

	tilesetPath := flag.String("tileset", "",
		`path to a schemas.tsj file`)
	flag.Parse()

	if *tilesetPath == "" {
		log.Fatal("--tileset can't be empty")
	}
	if len(flag.Args()) == 0 {
		log.Fatal("expected at least 1 positional argument")
	}

	tilesetData, err := os.ReadFile(*tilesetPath)
	if err != nil {
		log.Fatal(err)
	}

	tileset, err := tiled.UnmarshalTileset(tilesetData)
	if err != nil {
		log.Fatalf("[ERROR] decode tileset file: %v", err)
	}

	hasErrors := false
	for _, filename := range flag.Args() {
		err := checkFile(tileset, filename)
		if err != nil {
			hasErrors = true
			fmt.Fprintf(os.Stderr, "%q: %v\n", filename, err)
		}
	}

	if hasErrors {
		os.Exit(1)
	} else {
		fmt.Printf("[OK] all files are good (checked %d files)\n", len(flag.Args()))
	}
}

func checkFile(tileset *tiled.Tileset, filename string) error {
	levelData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return leveldata.ValidateLevelData(tileset, levelData)
}
