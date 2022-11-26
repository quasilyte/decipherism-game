package main

import (
	"fmt"

	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type generatorConfig struct {
	difficulty int
}

func generateLevel(r *gmath.Rand, tileset *tiled.Tileset, config generatorConfig) *schemaTemplate {
	g := levelGenerator{rand: r, config: config}
	return g.generate(tileset)
}

type levelGenerator struct {
	rand   *gmath.Rand
	config generatorConfig
	result *schemaTemplate
	tiles  [numSchemaRows][numSchemaCols]generatorTile
}

type generatorTile uint8

const (
	gtileUnset generatorTile = iota
	gtileOutOfBounds
	gtileElem
)

func (g *levelGenerator) generate(tileset *tiled.Tileset) *schemaTemplate {
	numAttempts := 1
	for {
		g.result = &schemaTemplate{tileset: tileset}
		ok := g.tryGenerate()
		if ok {
			break
		}
		g.tiles = [numSchemaRows][numSchemaCols]generatorTile{}
		numAttempts++
	}
	// Scale the positions: from coord to the actual pixel offsets.
	for i := range g.result.elems {
		elem := &g.result.elems[i]
		// elem.pos.X = elem.pos.X*tileset.TileWidth + tileset.TileWidth/2
		// elem.pos.Y = elem.pos.Y*tileset.TileHeight - tileset.TileHeight/2
		fmt.Println(elem.class, elem.pos, elem.rotation)
	}
	fmt.Printf("generation took %d attempts\n", numAttempts)
	return g.result
}

func (g *levelGenerator) tryGenerate() bool {
	tileset := g.result.tileset
	startCol := float64(g.rand.IntRange(2, numSchemaCols-2-1))
	startRow := float64(g.rand.IntRange(2, numSchemaRows-2-1))
	pos := gmath.Vec{
		X: startCol*tileset.TileWidth + tileset.TileWidth/2,
		Y: startRow*tileset.TileHeight + tileset.TileHeight/2,
	}
	// branchBudget := g.rand.IntRange(3+g.config.difficulty*2, 5+g.config.difficulty*2)
	branchBudget := 1
	g.deployElem(pos, "elem_input", 0)
	startShape := getElemShape("elem_input", pos, 0, nil)
	return g.buildBranch(pos, startShape, branchBudget)
}

func (g *levelGenerator) buildBranch(pos gmath.Vec, srcShape elemShape, budget int) bool {
	m := g.pickModule(budget)
	if m != nil {
		budget -= m.cost
	} else {
		m = outputModule
	}
	dirIndex := g.rand.IntRange(0, 3)
	for i := 0; i < len(allDirections); i++ {
		dir := allDirections[dirIndex]
		if exits, ok := g.tryDeployModule(m, srcShape, pos.MoveInDirection(96, dir), dir); ok {
			for _, e := range exits {
				shape := getElemShape(e.class, e.pos, e.rotation, e.extraData)
				if !g.buildBranch(e.pos, shape, budget) {
					return false
				}
			}
			return true
		}
		dirIndex++
		if dirIndex >= len(allDirections) {
			dirIndex = 0
		}
	}
	return false
}

func (g *levelGenerator) tryDeployModule(m *generatorModule, srcShape elemShape, pos gmath.Vec, dir gmath.Rad) ([]schemaTemplateElem, bool) {
	// Can this module be connected to a given source shape?
	entryElem := m.elems[0]
	entryElemPos := entryElem.pos.Add(pos)
	entryShape := getElemShape(entryElem.class, entryElemPos, entryElem.rotation, entryElem.extraData)
	if !srcShape.CanConnectTo(&entryShape) {
		return nil, false
	}
	// First check if we can deploy all module elements.
	for _, proto := range m.elems {
		elemPos := proto.pos.Add(pos)
		if g.readCell(elemPos) != gtileUnset {
			return nil, false
		}
	}
	// Now do the actual deployment.
	for _, proto := range m.elems {
		elemPos := proto.pos.Add(pos)
		g.deployElem(elemPos, proto.class, dir)
	}
	if m.numExits == 0 {
		return nil, true
	}
	return g.result.elems[len(g.result.elems)-m.numExits:], true
}

func (g *levelGenerator) pickModule(budget int) *generatorModule {
	index := g.rand.IntRange(0, len(generatorModules)-1)
	for i := 0; i < len(generatorModules); i++ {
		m := &generatorModules[index]
		if m.cost <= budget {
			return m
		}
		if index >= len(generatorModules) {
			index = 0
		}
	}
	return nil
}

func (g *levelGenerator) deployElem(pos gmath.Vec, class string, rotation gmath.Rad) {
	e := schemaTemplateElem{
		class:   class,
		pos:     pos,
		classID: -1,
	}
	switch class {
	case "pipe":
		e.rotation = rotation
	}
	g.result.elems = append(g.result.elems, e)
	g.markCell(pos, gtileElem)
}

func (g *levelGenerator) rowcolByPos(pos gmath.Vec) (int, int) {
	col := int(pos.X) / int(g.result.tileset.TileWidth)
	row := int(pos.Y) / int(g.result.tileset.TileHeight)
	return row, col
}

func (g *levelGenerator) markCell(pos gmath.Vec, t generatorTile) {
	row, col := g.rowcolByPos(pos)
	g.tiles[row][col] = t
}

func (g *levelGenerator) readCell(pos gmath.Vec) generatorTile {
	row, col := g.rowcolByPos(pos)
	if col >= numSchemaCols || col < 0 {
		return gtileOutOfBounds
	}
	if row >= numSchemaRows || row < 0 {
		return gtileOutOfBounds
	}
	return g.tiles[row][col]
}

var (
	dirRight = gmath.Rad(0)
	dirDown  = gmath.DegToRad(90)
	dirLeft  = gmath.DegToRad(180)
	dirUp    = gmath.DegToRad(270)
)

var allDirections = []gmath.Rad{
	dirRight,
	dirDown,
	dirLeft,
	dirUp,
}

type generatorModule struct {
	cost     int
	elems    []schemaTemplateElem
	numExits int
}

var outputModule = &generatorModule{
	numExits: 0,
	elems: []schemaTemplateElem{
		{class: "elem_output"},
	},
}

var generatorModules = []generatorModule{
	{
		cost:     1,
		numExits: 1,
		elems: []schemaTemplateElem{
			{class: "pipe"},
		},
	},
}
