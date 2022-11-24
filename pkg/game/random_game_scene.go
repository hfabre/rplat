package game

import (
	"fmt"
	"time"

	rl "github.com/chunqian/go-raylib/raylib"
)

const HookVerticalForce = 30
const HookHorizontalForce = 60

// Workaround to be able to call methods on a pointer on my interface
type RandGameSceneWrapper struct {
	rgs *RandomGameScene
}

func NewRandGameSceneWrapper() RandGameSceneWrapper {
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
	player          *Player
	level           Map
	inputManager    *InputManager
	elapsedSeconds  int
	ticker          *time.Ticker
	durationSeconds int
}

func NewRandomGameScene() *RandomGameScene {
	rgs := &RandomGameScene{}

	player := Player{
		pos:          rl.Vector2{32, 32},
		lastPos:      rl.Vector2{20, 20},
		velocity:     rl.Vector2{0, 0},
		lastVelocity: rl.Vector2{0, 0},
		size:         rl.Vector2{32, 64},
		canJump:      true,
		color:        rl.Red,
	}

	rgs.player = &player

	// Load level
	mc := NewMapConfiguration("./assets/map.json")
	tileset := NewTileset("./assets/tileset.png", mc.TileWidth, mc.TileHeight)
	rgs.level = NewMap(mc, tileset)

	im := NewInputManager()
	rgs.inputManager = &im

	rgs.durationSeconds = 30

	// TODO: Add an init function
	rgs.ticker = time.NewTicker(1 * time.Second)
	go rgs.StartsCounter()

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
			rgs.player.MoveRight()
		case "move_left":
			rgs.player.MoveLeft()
		case "jump":
			rgs.player.Jump()
		case "hook":
			rgs.player.Hook()
		case "stop_hook":
			rgs.player.StopHook()
		case "dash":
			rgs.player.Dash()
		case "portal":
			rgs.player.FirePortal(rgs.level.walls)
		default:
			// Unknown event
		}
	}
}

func (rgs *RandomGameScene) StartsCounter() {
	for _ = range rgs.ticker.C {
		rgs.elapsedSeconds++
	}
}

func (rgs *RandomGameScene) Update(deltaTime float32) {
	if rgs.elapsedSeconds >= rgs.durationSeconds {
		// Game finished
	}

	if !Pause {
		rgs.player.color = rl.Green
		rgs.player.lastPos = rgs.player.pos
		rgs.player.lastVelocity = rgs.player.velocity

		rgs.player.Update(deltaTime)

		// Resolve collisions
		for i := 0; i < len(rgs.level.walls); i++ {
			if isColliding(rgs.player.Rectangle(), rgs.level.walls[i]) {
				direction := collisionDirection(rgs.player.Rectangle(), rgs.level.walls[i])
				rgs.player.SolveCollision(rgs.level.walls[i], direction)
			}

			if rgs.player.hookLaunched {
				if isColliding(rgs.player.hook.Rectangle(), rgs.level.walls[i]) {
					direction := collisionDirection(rgs.player.hook.Rectangle(), rgs.level.walls[i])
					rgs.player.hook.SolveCollision(rgs.level.walls[i], direction)
				}
			}

			if rgs.player.portal.status == "ended" {
				if isColliding(rgs.player.Rectangle(), rgs.player.portal.EntryRectangle()) {
					rgs.player.Teleport(rgs.player.portal.exit_pos)
				}
			}
		}

		// if rl.Vector2Distance(rgs.player.lastPos, rgs.player.pos) > 10 {
		// 	Pause = true
		// }
	}
}

func (rgs RandomGameScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	rgs.level.Draw()

	currentStateLerp := LerpVec2(rgs.player.pos, factor)
	lastStateLerp := LerpVec2(rgs.player.pos, 1-factor)
	rl.DrawRectangleV(rl.Vector2{currentStateLerp.X + lastStateLerp.X, currentStateLerp.Y + lastStateLerp.Y}, rgs.player.size, rgs.player.color)

	if rgs.player.hookLaunched {
		currentStateLerp = LerpVec2(rgs.player.hook.pos, factor)
		lastStateLerp = LerpVec2(rgs.player.hook.pos, 1-factor)
		rl.DrawRectangleV(rl.Vector2{currentStateLerp.X + lastStateLerp.X, currentStateLerp.Y + lastStateLerp.Y}, rgs.player.hook.size, rgs.player.hook.color)

		rl.DrawLineEx(rgs.player.pos, rgs.player.hook.pos, 5, rl.Black)
	}

	switch rgs.player.portal.status {
	case "triggered":
		rl.DrawRectangleV(rgs.player.portal.entry_pos, rl.Vector2{PortalWidth, PortalHeight}, rl.Blue)
	case "ended":
		rl.DrawRectangleV(rgs.player.portal.entry_pos, rl.Vector2{PortalWidth, PortalHeight}, rl.Blue)
		rl.DrawRectangleV(rgs.player.portal.exit_pos, rl.Vector2{PortalWidth, PortalHeight}, rl.Brown)
	}

	timeText := fmt.Sprintf("Elapsed time: %v", rgs.elapsedSeconds)
	rl.DrawText(timeText, 500, 20, 40, rl.Black)

	if Debug {
		if Pause {
			rl.DrawRectangleV(rl.Vector2{rgs.player.lastPos.X, rgs.player.lastPos.Y}, rgs.player.size, rl.Gray)
		}

		posText := fmt.Sprintf("Position: %v - %v", rgs.player.pos.X, rgs.player.pos.Y)
		lastPosText := fmt.Sprintf("Last position: %v - %v", rgs.player.lastPos.X, rgs.player.lastPos.Y)
		velText := fmt.Sprintf("Velocity: %v - %v", rgs.player.velocity.X, rgs.player.velocity.Y)
		lastVelText := fmt.Sprintf("Last Velocity: %v - %v", rgs.player.lastVelocity.X, rgs.player.lastVelocity.Y)
		distText := fmt.Sprintf("Distance: %v", rl.Vector2Distance(rgs.player.lastPos, rgs.player.pos))

		rl.DrawFPS(10, 10)
		rl.DrawText(posText, 10, 50, 20, rl.Black)
		rl.DrawText(lastPosText, 10, 70, 20, rl.Black)
		rl.DrawText(velText, 10, 90, 20, rl.Black)
		rl.DrawText(lastVelText, 10, 110, 20, rl.Black)
		rl.DrawText(distText, 10, 130, 20, rl.Black)
	}
}
