package game

import (
	"fmt"
	rl "github.com/chunqian/go-raylib/raylib"
)

const FPS = 500

// In pixel per second
const PlayerSpeed = 100
const PlayerJumpSpeed = 500
const HookSpeed = 1800
const HookVerticalForce = 30
const HookHorizontalForce = 60

const Friction = 0.80
const Gravity = 10

// Workaround to be able to call methods on a pointer on my interface
type RandGameSceneWrapper struct {
	rgs RandomGameScene
}

func NewRandGameSceneWrapper () RandGameSceneWrapper {
	return RandGameSceneWrapper{rgs: NewRandomGameScene()}
}

// Implement Scene interface
func (rgsw RandGameSceneWrapper) Init() {

}

func (rgsw RandGameSceneWrapper) UpdateInputs() {
	rgsw.rgs.UpdateInputs()
}

func (rgsw RandGameSceneWrapper) ClearInputs() {
	rgsw.rgs.ClearInputs()
}

func (rgsw RandGameSceneWrapper) HandleEvents() {
	rgsw.rgs.HandleEvents()
}

func (rgsw RandGameSceneWrapper) Update(dt float32) {
	rgsw.rgs.Update(dt)
}

func (rgsw RandGameSceneWrapper) Draw(factor float64) {
	rgsw.rgs.Draw(factor)
}

func (rgsw RandGameSceneWrapper) End() {

}

type RandomGameScene struct {
	player *Player
	level Map
	inputManager *InputManager
}

func NewRandomGameScene() RandomGameScene {
	rgs := RandomGameScene{}

	player := Player{
		pos:      rl.Vector2{32, 32},
		lastPos:  rl.Vector2{20, 20},
		velocity: rl.Vector2{0, 0},
		lastVelocity: rl.Vector2{0, 0},
		size:     rl.Vector2{32, 64},
		canJump: true,
		color: rl.Red,
	}

	rgs.player = &player

	// Load level
	mc := NewMapConfiguration("./assets/map.json")
	tileset := NewTileset("./assets/tileset.png", mc.TileWidth, mc.TileHeight)
	rgs.level = NewMap(mc, tileset)

	im := NewInputManager()
	rgs.inputManager = &im
	return rgs
}

func (rgs *RandomGameScene) UpdateInputs() {
	rgs.inputManager.Update()
}

func (rgs *RandomGameScene) ClearInputs() {
	rgs.inputManager.Clear()
}

func (rgs *RandomGameScene) HandleEvents() {
	for i := 0; i < len(rgs.inputManager.events); i++ {
		switch e := rgs.inputManager.events[i]; e {
		case "pause":
			Pause = !Pause
		case "move_right":
			rgs.player.velocity.X += PlayerSpeed
		case "move_left":
			rgs.player.velocity.X -= PlayerSpeed
		case "jump":
			if rgs.player.canJump {
				rgs.player.canJump = false
				rgs.player.velocity.Y -= PlayerJumpSpeed
			}
		case "hook":
			if !rgs.player.hookLaunched {
				rgs.player.hook = NewHook(*rgs.player)
				rgs.player.hookLaunched = true
			}
		case "stop_hook":
			rgs.player.hookLaunched = false
		default:
			// Unknown event
		}
	}
}

func (rgs *RandomGameScene) Update(deltaTime float32) {
	if !Pause {
		rgs.player.color = rl.Green
		rgs.player.lastPos = rgs.player.pos
		rgs.player.lastVelocity = rgs.player.velocity


		rgs.player.Update(deltaTime)

		// Resolve collisions
		for i := 0; i < len(rgs.level.walls); i++ {
			if rl.CheckCollisionRecs(rgs.player.Rectangle(), rgs.level.walls[i]) {
				rgs.player.SolveCollision(rgs.level.walls[i])
			}

			if rgs.player.hookLaunched {
				if rl.CheckCollisionRecs(rgs.player.hook.Rectangle(), rgs.level.walls[i]) {
					rgs.player.hook.SolveCollision(rgs.level.walls[i])
				}
			}
		}
	}
}

func (rgs RandomGameScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	rgs.level.Draw()

	currentStateLerp := LerpVec2(rgs.player.pos, factor)
	lastStateLerp := LerpVec2(rgs.player.pos, 1 - factor)
	rl.DrawRectangleV(rl.Vector2{currentStateLerp.X + lastStateLerp.X, currentStateLerp.Y + lastStateLerp.Y}, rgs.player.size, rgs.player.color)

	if rgs.player.hookLaunched {
		currentStateLerp = LerpVec2(rgs.player.hook.pos, factor)
		lastStateLerp = LerpVec2(rgs.player.hook.pos, 1 - factor)
		rl.DrawRectangleV(rl.Vector2{currentStateLerp.X + lastStateLerp.X, currentStateLerp.Y + lastStateLerp.Y}, rgs.player.hook.size, rgs.player.hook.color)

		rl.DrawLineEx(rgs.player.pos, rgs.player.hook.pos, 5, rl.Black)
	}

	if Debug {
		if Pause {
			rl.DrawRectangleV(rl.Vector2{rgs.player.collision.X, rgs.player.collision.Y}, rl.Vector2{10, 10}, rl.Blue)
		}

		posText := fmt.Sprintf("Position: %v - %v", rgs.player.pos.X, rgs.player.pos.Y)
		lastPosText := fmt.Sprintf("Last position: %v - %v", rgs.player.lastPos.X, rgs.player.lastPos.Y)
		velText := fmt.Sprintf("Velocity: %v - %v", rgs.player.velocity.X, rgs.player.velocity.Y)
		lastVelText := fmt.Sprintf("Last Velocity: %v - %v", rgs.player.lastVelocity.X, rgs.player.lastVelocity.Y)
		colText := fmt.Sprintf("Last collision: %v", rgs.player.collision)

		rl.DrawFPS(10, 10)
		rl.DrawText(posText, 10, 50, 20, rl.Black)
		rl.DrawText(lastPosText, 10, 70, 20, rl.Black)
		rl.DrawText(velText, 10, 90, 20, rl.Black)
		rl.DrawText(lastVelText, 10, 110, 20, rl.Black)
		rl.DrawText(colText, 10, 130, 20, rl.Black)
	}
}
