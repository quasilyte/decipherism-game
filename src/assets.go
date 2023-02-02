package main

import (
	"strings"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
)

func prepareAssets(ctx *ge.Context) {
	theStoryModeMap.levels = make(map[string]storyModeLevel)
	resourceID := rawLastID + 1
	resourceID = loadLevelsData(ctx, resourceID, "levels/story")
	resourceID = loadLevelsData(ctx, resourceID, "levels/bonus")
}

func loadLevelsData(ctx *ge.Context, idSeq resource.RawID, dirname string) resource.RawID {
	levels, err := gameAssets.ReadDir("_assets/" + dirname)
	if err != nil {
		panic(err)
	}
	for _, f := range levels {
		shortName := strings.TrimSuffix(f.Name(), ".json")
		theStoryModeMap.levels[shortName] = storyModeLevel{
			name: shortName,
			id:   idSeq,
		}
		ctx.Loader.RawRegistry.Set(idSeq, resource.RawInfo{
			Path: dirname + "/" + f.Name(),
		})
		ctx.Loader.LoadRaw(idSeq)
		idSeq++
	}
	return idSeq
}
