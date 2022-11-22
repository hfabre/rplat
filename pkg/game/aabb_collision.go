package game

import (
	"math"

	rl "github.com/chunqian/go-raylib/raylib"
)

func isColliding(main_body, other_body rl.Rectangle) bool {
	var p1x = math.Max(float64(main_body.X), float64(other_body.X))
	var p1y = math.Max(float64(main_body.Y), float64(other_body.Y))
	var p2x = math.Min(float64(main_body.X+main_body.Width), float64(other_body.X+other_body.Width))
	var p2y = math.Min(float64(main_body.Y+main_body.Height), float64(other_body.Y+other_body.Height))

	if (p2x-p1x) > 0 && (p2y-p1y) > 0 {
		return true
	} else {
		return false
	}
}

func collisionDirection(main_body, other_body rl.Rectangle) string {
	left := (main_body.X + main_body.Width) - other_body.X
	right := (other_body.X + other_body.Width) - main_body.X
	bottom := (main_body.Y + main_body.Height) - other_body.Y
	top := (other_body.Y + other_body.Height) - main_body.Y

	if right < left && right < top && right < bottom {
		return "right"
	} else if left < top && left < bottom {
		return "left"
	} else if top < bottom {
		return "top"
	} else {
		return "bottom"
	}
}
