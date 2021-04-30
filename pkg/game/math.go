package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
	"math"
)

func angleFromVectors (v1, v2 rl.Vector2) float64 {
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
	return rl.Vector2{X: vec.X * float32(factor), Y: vec.Y * float32(factor) }
}
