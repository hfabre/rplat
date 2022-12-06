package game

import (
	"time"

	rl "github.com/chunqian/go-raylib/raylib"
)

const Friction = 0.80
const Gravity = 10

const PlayerSpeed = 100
const PlayerJumpSpeed = 550
const DashForce = 8
const DashCooldown = 500
const PortalCooldown = 500

type Player struct {
	pos, lastPos, velocity, lastVelocity, hookVelocity, size rl.Vector2
	canJump, hookLaunched                                    bool
	color                                                    rl.Color
	hook                                                     Hook
	last_dash_time, last_portal_time                         int64
	portal                                                   Portal
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

func (p *Player) MoveRight() {
	p.velocity.X += PlayerSpeed
}

func (p *Player) MoveLeft() {
	p.velocity.X -= PlayerSpeed
}

func (p *Player) Jump() {
	if p.canJump {
		p.canJump = false
		p.velocity.Y -= PlayerJumpSpeed
	}
}

func (p *Player) Dash() {
	current_time := time.Now().UnixNano() / int64(time.Millisecond)

	if current_time-p.last_dash_time > DashCooldown {
		p.last_dash_time = current_time
		p.velocity.X = p.velocity.X * DashForce
	}
}

func (p *Player) Hook() {
	if !p.hookLaunched {
		p.hook = NewHook(*p)
		p.hookLaunched = true
	}
}

func (p *Player) StopHook() {
	p.hookLaunched = false
}

func (p *Player) FirePortal(walls []rl.Rectangle) {
	current_time := time.Now().UnixNano() / int64(time.Millisecond)

	if current_time-p.last_portal_time > PortalCooldown {
		p.last_portal_time = current_time
		portal_box := p.Rectangle()
		dir := DirectionVectorFromVectors(p.pos, rl.GetMousePosition())
		velocity := rl.Vector2{X: dir.X * 10, Y: dir.Y * 10}

		// TODO: impl Ray -> AABB collision
		for j := 0; j < 10000; j++ {
			for i := 0; i < len(walls); i++ {
				if isColliding(portal_box, walls[i]) {
					direction := collisionDirection(portal_box, walls[i])
					SolvePortalCollision(&portal_box, walls[i], direction)
					p.portal.Trigger(rl.Vector2{X: portal_box.X, Y: portal_box.Y})
					return
				}
			}

			portal_box.X += velocity.X
			portal_box.Y += velocity.Y
		}
	}
}

func (p *Player) Teleport(pos rl.Vector2) {
	p.pos.X = pos.X
	p.pos.Y = pos.Y
}

// Note: Hook physics is heavily inspired by Teeworlds, see:
// https://github.com/teeworlds/teeworlds/blob/b0c4c7002b28ee195934281e524f163f7ed30c59/src/game/gamecore.cpp#L263
func (p *Player) Update(deltaTime float32) {
	if p.hookLaunched {
		if p.hook.hooked {
			dir := DirectionVectorFromVectors(p.pos, p.hook.pos)
			p.hookVelocity.X = dir.X * HookHorizontalForce
			p.hookVelocity.Y = dir.Y * HookVerticalForce

			// The hook as more power to drag you up then down. This makes it easier to get on top of a platform
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

func (p *Player) checkAndHandleCollisions(walls []rl.Rectangle) {
	for i := 0; i < len(walls); i++ {
		if isColliding(p.Rectangle(), walls[i]) {
			direction := collisionDirection(p.Rectangle(), walls[i])
			p.SolveCollision(walls[i], direction)
		}

		if p.hookLaunched {
			if isColliding(p.hook.Rectangle(), walls[i]) {
				direction := collisionDirection(p.hook.Rectangle(), walls[i])
				p.hook.SolveCollision(walls[i], direction)
			}
		}

		if p.portal.status == "ended" {
			if isColliding(p.Rectangle(), p.portal.EntryRectangle()) {
				p.StopHook()
				p.Teleport(p.portal.exit_pos)
			}
		}
	}
}

func (p *Player) SolveCollision(wall rl.Rectangle, direction string) {
	p.color = rl.Red

	switch direction {
	case "bottom":
		p.canJump = true
		p.pos.Y = wall.Y - p.size.Y
		p.velocity.Y = 0
	case "right":
		p.pos.X = wall.X + wall.Width
		p.velocity.X = 0
	case "left":
		p.pos.X = wall.X - p.size.X
		p.velocity.X = 0
	case "top":
		p.pos.Y = wall.Y + wall.Height
		p.velocity.Y = 0
	}
}
