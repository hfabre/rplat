package game

import rl "github.com/chunqian/go-raylib/raylib"

type Hook struct {
	pos, lastPos, velocity, size rl.Vector2
	hooked bool
	color rl.Color
}

func NewHook(player Player) Hook {
	dir := DirectionVectorFromVectors(player.pos, rl.GetMousePosition())

	return Hook{
		pos:      player.pos,
		lastPos:  player.pos,
		velocity: rl.Vector2{dir.X * HookSpeed, dir.Y * HookSpeed},
		size:     rl.Vector2{32, 32},
		hooked: false,
		color: rl.Orange,
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

func (h *Hook) SolveCollision (wall rl.Rectangle) {
	h.color = rl.Purple
	collision := rl.GetCollisionRec(h.Rectangle(), wall)

	// If collision is too small just ignore it
	// It avoid having a bug where shallow axis is not the expected one on very small values
	if collision.Width < 1 && collision.Height < 1 {
		return
	}

	if collision.Width > collision.Height {

		// If the hook is going down
		if h.velocity.Y > 0 {
			h.pos.Y = collision.Y - h.size.Y
		} else {
			h.pos.Y = collision.Y + collision.Height
		}
	} else {
		if h.velocity.X > 0 {
			h.pos.X = collision.X - h.size.X
		} else {
			h.pos.X = collision.X + collision.Width
		}
	}

	h.velocity.X = 0
	h.velocity.Y = 0
	h.hooked = true
}