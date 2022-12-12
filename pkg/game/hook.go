package game

import rl "github.com/chunqian/go-raylib/raylib"

const HookSpeed = 1800
const HookVerticalForce = 30
const HookHorizontalForce = 60

type Hook struct {
	pos, lastPos, velocity, size rl.Vector2
	hooked                       bool
	color                        rl.Color
}

func NewHook(player Player) Hook {
	dir := DirectionVectorFromVectors(player.pos, rl.GetMousePosition())

	return Hook{
		pos:      player.pos,
		lastPos:  player.pos,
		velocity: rl.Vector2{X: dir.X * HookSpeed, Y: dir.Y * HookSpeed},
		size:     rl.Vector2{X: 32, Y: 32},
		hooked:   false,
		color:    rl.Orange,
	}
}

func (h Hook) Rectangle() rl.Rectangle {
	// TODO: Avoid creating new rectangle each time if this becomes a performance bottleneck
	return rl.Rectangle{
		X:      h.pos.X,
		Y:      h.pos.Y,
		Width:  h.size.X,
		Height: h.size.Y,
	}
}

func (h *Hook) SolveCollision(wall rl.Rectangle, direction string) {
	switch direction {
	case "bottom":
		h.pos.Y = wall.Y - h.size.Y
		h.velocity.Y = 0
	case "right":
		h.pos.X = wall.X + wall.Width
		h.velocity.X = 0
	case "left":
		h.pos.X = wall.X - h.size.X
		h.velocity.X = 0
	case "top":
		h.pos.Y = wall.Y + wall.Height
		h.velocity.Y = 0
	}

	h.velocity.X = 0
	h.velocity.Y = 0
	h.hooked = true
}
