package main

import (
	"hash/fnv"

	"github.com/quasilyte/decipherism-game/leveldata"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/tiled"
)

func loadLevelTemplate(scene *ge.Scene, levelData []byte) (*leveldata.SchemaTemplate, error) {
	tileset, err := tiled.UnmarshalTileset(scene.LoadRaw(RawComponentSchemaTilesetJSON).Data)
	if err != nil {
		panic(err)
	}
	return leveldata.LoadLevelTemplate(tileset, levelData)
}

func fnvhash(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

func volumeMultiplier(level int) float64 {
	switch level {
	case 1:
		return 0.05
	case 2:
		return 0.10
	case 3:
		return 0.3
	case 4:
		return 0.55
	case 5:
		return 0.8
	default:
		return 0
	}
}
