package game

import (
	"encoding/json"
	"fmt"
	rl "github.com/chunqian/go-raylib/raylib"
	"io/ioutil"
)

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Tile struct {
	Index      int        `json:"index"`
	Properties []Property `json:"properties"`
}

//
//  Tileset
//

type Tileset struct {
	tiles []rl.Texture2D
}

func NewTileset(path string, tileWidth, tileHeight int) Tileset {
	var tiles []rl.Texture2D
	tileset := rl.LoadImage(path)
	horizontalTileCount := int(tileset.Width) / tileWidth
	verticalTileCount := int(tileset.Height) / tileHeight

	for y := 0; y < verticalTileCount; y++ {
		for x := 0; x < horizontalTileCount; x++ {
			rect := rl.Rectangle{X: float32(x * tileWidth), Y: float32(y * tileHeight), Width: float32(tileWidth), Height: float32(tileHeight)}
			image := rl.ImageFromImage(tileset, rect)
			tiles = append(tiles, rl.LoadTextureFromImage(image))
		}
	}

	return Tileset{tiles: tiles}
}

func (ts Tileset) Unload() {
	for _, tile := range ts.tiles {
		rl.UnloadTexture(tile)
	}
}

//
//  MapConfiguration
//

type MapConfiguration struct {
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	TileHeight int      `json:"tileHeight"`
	TileWidth  int      `json:"tileWidth"`
	Board      [][]Tile `json:"tiles"`
}

func NewMapConfiguration(path string) MapConfiguration {
	var mc MapConfiguration

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(data, &mc)
	if err != nil {
		fmt.Println("error:", err)
	}

	return mc
}

//
//  Map
//

type Map struct {
	mc MapConfiguration
	ts Tileset

	width      int
	height     int
	tileWidth  int
	tileHeight int
	board      [][]Tile
	walls      []rl.Rectangle
}

func NewMap(mc MapConfiguration, ts Tileset) Map {
	var walls []rl.Rectangle

	for y := 0; y < len(mc.Board); y++ {
		for x := 0; x < len(mc.Board[y]); x++ {
			tileIndex := mc.Board[y][x].Index

			// TODO: Handle properties
			if tileIndex >= 0 {
				walls = append(walls, rl.Rectangle{
					X:      float32(x * mc.TileWidth),
					Y:      float32(y * mc.TileHeight),
					Width:  float32(mc.TileWidth),
					Height: float32(mc.TileHeight),
				})
			}
		}
	}

	m := Map{
		mc:         mc,
		ts:         ts,
		width:      mc.Width,
		height:     mc.Height,
		tileWidth:  mc.TileWidth,
		tileHeight: mc.TileHeight,
		board:      mc.Board,
		walls:      walls,
	}

	return m
}

func (m Map) Draw() {
	for y := 0; y < len(m.board); y++ {
		for x := 0; x < len(m.board[y]); x++ {
			tileIndex := m.board[y][x].Index

			if tileIndex >= 0 {
				rl.DrawTexture(m.ts.tiles[tileIndex], int32(x*m.tileWidth), int32(y*m.tileHeight), rl.White)
			}
		}
	}
}
