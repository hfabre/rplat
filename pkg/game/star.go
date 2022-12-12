package game

import rl "github.com/chunqian/go-raylib/raylib"

const StarWidth = 32
const StarHeight = 32

type Star struct {
	pos rl.Vector2
}

func (s Star) Rectangle() rl.Rectangle {
	return rl.Rectangle{
		X:      s.pos.X,
		Y:      s.pos.Y,
		Width:  StarWidth,
		Height: StarHeight,
	}
}

func (s Star) Draw() {
	rl.DrawRectangleRec(s.Rectangle(), rl.Yellow)
}
