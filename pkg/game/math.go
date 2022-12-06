package game

import (
	"math"

	rl "github.com/chunqian/go-raylib/raylib"
)

func angleFromVectors(v1, v2 rl.Vector2) float64 {
	x := v2.X - v1.X
	y := v2.Y - v1.Y

	return math.Atan2(float64(y), float64(x))
}

func DirectionVectorFromAngle(angle float64) rl.Vector2 {
	return rl.Vector2{
		X: float32(math.Cos(angle)),
		Y: float32(math.Sin(angle)),
	}
}

func DirectionVectorFromVectors(v1, v2 rl.Vector2) rl.Vector2 {
	return DirectionVectorFromAngle(angleFromVectors(v1, v2))
}

func LerpVec2(vec rl.Vector2, factor float64) rl.Vector2 {
	return rl.Vector2{X: vec.X * float32(factor), Y: vec.Y * float32(factor)}
}

// See https://github.com/MonoGame/MonoGame/blob/2911bbb3bcc412aedc00c109f3a96bb2480b255f/MonoGame.Framework/Ray.cs#L97
// Does not works as expected.
func RayAABBCollision(origin, direction rl.Vector2, aabb rl.Rectangle) (bool, float32) {
	epsilon := 0.000001
	box_min := rl.Vector2{X: aabb.X, Y: aabb.Y}
	box_max := rl.Vector2{X: aabb.X + aabb.Width, Y: aabb.Y + aabb.Height}
	tMin := float32(math.Inf(-1))
	tMax := float32(math.Inf(1))

	if math.Abs(float64(direction.X)) < epsilon {
		if origin.X < box_min.X || origin.X > box_max.X {
			return false, 0
		}
	} else {
		tMin = (box_min.X - origin.X) / direction.X
		tMax = (box_max.X - origin.X) / direction.X

		if tMin > tMax {
			temp := tMin
			tMin = tMax
			tMax = temp
		}
	}

	if math.Abs(float64(direction.Y)) < epsilon {
		if origin.Y < box_min.Y || origin.Y > box_max.Y {
			return false, 0
		}
	} else {
		var tMinY = (box_min.Y - origin.Y) / direction.Y
		var tMaxY = (box_max.Y - origin.Y) / direction.Y

		if tMinY > tMaxY {
			temp := tMinY
			tMinY = tMaxY
			tMaxY = temp
		}

		if (tMin != float32(math.Inf(-1)) && tMin > tMaxY) || (tMax != float32(math.Inf(1)) && tMinY > tMax) {
			return false, 0
		}

		if tMin == float32(math.Inf(-1)) || tMinY > tMin {
			tMin = tMinY
		}

		if tMax == float32(math.Inf(1)) || tMaxY < tMax {
			tMax = tMaxY
		}
	}

	// having a positive tMax and a negative tMin means the ray is inside the box
	// we expect the intesection distance to be 0 in that case
	if (tMin != float32(math.Inf(-1)) && tMin < 0) && tMax > 0 {
		return true, 0
	}

	// a negative tMin means that the intersection point is behind the ray's origin
	// we discard these as not hitting the AABB
	if tMin < 0 {
		return false, 0
	}

	return true, tMin
}
