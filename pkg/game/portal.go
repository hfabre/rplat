package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
)

const PortalWidth = 32
const PortalHeight = 32

type Portal struct {
	entry_pos, exit_pos rl.Vector2
	status              string
}

func (p *Portal) Trigger(pos rl.Vector2) {
	if p.status == "triggered" {
		if !isColliding(p.EntryRectangle(), rl.Rectangle{X: pos.X, Y: pos.Y, Width: PortalWidth, Height: PortalHeight}) {
			p.exit_pos = pos
			p.status = "ended"
		}
	} else {
		p.entry_pos = pos
		p.status = "triggered"
	}
}

func (p Portal) EntryRectangle() rl.Rectangle {
	return rl.Rectangle{
		X:      p.entry_pos.X,
		Y:      p.entry_pos.Y,
		Width:  PortalWidth,
		Height: PortalHeight,
	}
}

func SolvePortalCollision(portal_box *rl.Rectangle, wall rl.Rectangle, direction string) {
	switch direction {
	case "bottom":
		portal_box.Y = wall.Y - PortalHeight
	case "right":
		portal_box.X = wall.X + wall.Width
	case "left":
		portal_box.X = wall.X - PortalWidth
	case "top":
		portal_box.Y = wall.Y + wall.Height
	}
}
