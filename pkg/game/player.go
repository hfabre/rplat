package game

import rl "github.com/chunqian/go-raylib/raylib"

type Player struct {
	pos, lastPos, velocity, lastVelocity, hookVelocity, size rl.Vector2
	canJump, hookLaunched bool
	collision rl.Rectangle
	color rl.Color
	hook Hook
}

func (p Player) Rectangle() rl.Rectangle {
	// TODO: Avoid creating new rectangle each time if this becomes a performance bottleneck
	return rl.Rectangle{
		X:      p.pos.X,
		Y:      p.pos.Y,
		Width:  p.size.X,
		Height: p.size.Y,
	}
}


// Note: Hook physics is heavily inspired by Teeworlds, see:
// https://github.com/teeworlds/teeworlds/blob/b0c4c7002b28ee195934281e524f163f7ed30c59/src/game/gamecore.cpp#L263
func (p *Player) Update(deltaTime float32) {
	if p.hookLaunched {
		if p.hook.hooked {
			dir := DirectionVectorFromVectors(p.pos, p.hook.pos)
			p.hookVelocity.X = dir.X * HookHorizontalForce
			p.hookVelocity.Y = dir.Y * HookVerticalForce

			// The hook as more power to drag you up then down. This makes it easier to get on top of an platform
			if p.hookVelocity.Y > 0 {
				p.hookVelocity.Y *= 0.3
			}

			// The hook will boost it's power if the player wants to move on that direction.
			// Otherwise it will slow down everything a bit
			if p.hookVelocity.X < 0 && p.velocity.X < 0 || p.hookVelocity.X > 0 && p.velocity.X > 0 {
				p.hookVelocity.X *= 0.95
			} else {
				p.hookVelocity.X *= 0.75
			}

			// Apply hook physics
			p.velocity.X += p.hookVelocity.X
			p.velocity.Y += p.hookVelocity.Y

		} else {
			p.hook.pos.X += p.hook.velocity.X * deltaTime
			p.hook.pos.Y += p.hook.velocity.Y * deltaTime
		}
	}

	// Run natural forces
	p.velocity.X *= Friction
	p.velocity.Y += Gravity

	// Apply velocity
	p.pos.X += p.velocity.X * deltaTime
	p.pos.Y += p.velocity.Y * deltaTime
}

func (p *Player) SolveCollision(wall rl.Rectangle) {
	p.color = rl.Red
	collision := rl.GetCollisionRec(p.Rectangle(), wall)

	// If collision is too small just ignore it
	// It avoid having a bug where shallow axis is not the expected one on very small values
	if collision.Width < 1 && collision.Height < 1 {
		return
	}

	p.collision = collision

	// Perform collision resolution on shallow axis
	// Shallow axis is the less penetrated, so here Y axis
	if collision.Width > collision.Height {

		// If the player is going down
		if p.velocity.Y > 0 {
			p.canJump = true
			p.pos.Y = collision.Y - p.size.Y
		} else {
			p.pos.Y = collision.Y + collision.Height
		}

		p.velocity.Y = 0
	} else {
		// If the player is going right
		if p.velocity.X > 0 {
			p.pos.X = collision.X - p.size.X
		} else {
			p.pos.X = collision.X + collision.Width
		}

		p.velocity.X = 0
	}
}