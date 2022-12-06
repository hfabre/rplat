package game

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/chunqian/go-raylib/raylib"
)

const HookVerticalForce = 30
const HookHorizontalForce = 60

// Workaround to be able to call methods on a pointer on my interface
type RandGameSceneWrapper struct {
	rgs *RandomGameScene
}

func NewRandGameSceneWrapper(sm *SceneManager) RandGameSceneWrapper {
	return RandGameSceneWrapper{rgs: NewRandomGameScene(sm)}
}

// Implement Scene interface
func (rgsw RandGameSceneWrapper) Init() {
	rgsw.rgs.Init()
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
	rgsw.rgs.End()
}

func (rgsw RandGameSceneWrapper) ShouldExit() bool {
	return rgsw.rgs.ShouldExit()
}

type RandomGameScene struct {
	player          *Player
	level           Map
	inputManager    *InputManager
	elapsedSeconds  int
	ticker          *time.Ticker
	durationSeconds int
	stars           []Star
	score           int
	sceneManager    *SceneManager
	gameEnded       bool
}

func NewRandomGameScene(sm *SceneManager) *RandomGameScene {
	rgs := &RandomGameScene{}

	// Load level
	mc := NewMapConfiguration("./assets/map.json")
	tileset := NewTileset("./assets/tileset.png", mc.TileWidth, mc.TileHeight)
	im := NewInputManager()

	rgs.level = NewMap(mc, tileset)
	rgs.inputManager = &im
	rgs.durationSeconds = 30
	rgs.sceneManager = sm

	return rgs
}

func (rgs *RandomGameScene) Init() {
	player := Player{
		pos:          rl.Vector2{X: 32, Y: 32},
		lastPos:      rl.Vector2{X: 20, Y: 20},
		velocity:     rl.Vector2{X: 0, Y: 0},
		lastVelocity: rl.Vector2{X: 0, Y: 0},
		size:         rl.Vector2{X: 32, Y: 64},
		canJump:      true,
		color:        rl.Red,
	}

	rgs.player = &player
	rgs.gameEnded = false
	rgs.ticker = time.NewTicker(1 * time.Second)
	rgs.elapsedSeconds = 0

	for i := 0; i < 20; i++ {
		for !rgs.SpawnStar() {
		}
	}

	go rgs.StartsCounter()
}

func (rgs *RandomGameScene) End() {
	rgs.stars = nil
}

func (rgs *RandomGameScene) SpawnStar() bool {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	x := r.Intn(ScreenWidth)
	y := r.Intn(ScreenHeight)

	star := Star{rl.Vector2{X: float32(x), Y: float32(y)}}

	for _, wall := range rgs.level.walls {
		if isColliding(star.Rectangle(), wall) {
			return false
		}
	}

	for _, exStar := range rgs.stars {
		if isColliding(star.Rectangle(), exStar.Rectangle()) {
			return false
		}
	}

	rgs.stars = append(rgs.stars, star)
	return true
}

func (rgs *RandomGameScene) UpdateInputs() {
	rgs.inputManager.Update()
}

func (rgs *RandomGameScene) ClearInputs() {
	rgs.inputManager.Clear()
}

func (rgs *RandomGameScene) HandleEvents() {
	for i := 0; i < len(rgs.inputManager.events); i++ {
		if rgs.gameEnded && rgs.inputManager.events[i] == "validate" {
			rgs.sceneManager.SwapScene("main_menu")
		} else {
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
}

func (rgs *RandomGameScene) StartsCounter() {
	for range rgs.ticker.C {
		if !Pause {
			rgs.elapsedSeconds++

			if rgs.elapsedSeconds >= rgs.durationSeconds {
				rgs.EndGame(false)
				return
			}
		}
	}
}

func (rgs *RandomGameScene) EndGame(complete bool) {
	if complete {
		multiplier := rgs.durationSeconds - rgs.elapsedSeconds
		rgs.score *= multiplier
	}

	rgs.gameEnded = true
	rgs.ticker.Stop()
}

func (rgs RandomGameScene) ShouldExit() bool {
	return false
}

func (rgs *RandomGameScene) Update(deltaTime float32) {
	if rgs.gameEnded {
		return
	}

	if !Pause {
		rgs.player.color = rl.Green
		rgs.player.lastPos = rgs.player.pos
		rgs.player.lastVelocity = rgs.player.velocity

		rgs.player.Update(deltaTime)
		rgs.player.checkAndHandleCollisions(rgs.level.walls)

		starsToRemove := []int{}
		for i, star := range rgs.stars {
			if isColliding(rgs.player.Rectangle(), star.Rectangle()) {
				rgs.score += 10
				starsToRemove = append(starsToRemove, i)
			}
		}

		for _, j := range starsToRemove {
			rgs.stars = RemoveIndexStar(rgs.stars, j)
		}

		if len(rgs.stars) == 0 {
			rgs.EndGame(true)
		}
	}
}

func (rgs RandomGameScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	rgs.level.Draw()

	currentStateLerp := LerpVec2(rgs.player.pos, factor)
	lastStateLerp := LerpVec2(rgs.player.pos, 1-factor)
	rl.DrawRectangleV(rl.Vector2{X: currentStateLerp.X + lastStateLerp.X, Y: currentStateLerp.Y + lastStateLerp.Y}, rgs.player.size, rgs.player.color)

	if rgs.player.hookLaunched {
		currentStateLerp = LerpVec2(rgs.player.hook.pos, factor)
		lastStateLerp = LerpVec2(rgs.player.hook.pos, 1-factor)
		rl.DrawRectangleV(rl.Vector2{X: currentStateLerp.X + lastStateLerp.X, Y: currentStateLerp.Y + lastStateLerp.Y}, rgs.player.hook.size, rgs.player.hook.color)

		rl.DrawLineEx(rgs.player.pos, rgs.player.hook.pos, 5, rl.Black)
	}

	switch rgs.player.portal.status {
	case "triggered":
		rl.DrawRectangleV(rgs.player.portal.entry_pos, rl.Vector2{X: PortalWidth, Y: PortalHeight}, rl.Blue)
	case "ended":
		rl.DrawRectangleV(rgs.player.portal.entry_pos, rl.Vector2{X: PortalWidth, Y: PortalHeight}, rl.Blue)
		rl.DrawRectangleV(rgs.player.portal.exit_pos, rl.Vector2{X: PortalWidth, Y: PortalHeight}, rl.Brown)
	}

	for _, star := range rgs.stars {
		rl.DrawRectangleRec(star.Rectangle(), rl.Yellow)
	}

	timeText := fmt.Sprintf("Elapsed time: %v", rgs.elapsedSeconds)
	rl.DrawText(timeText, 500, 20, 40, rl.Black)

	if rgs.gameEnded {
		timeText := fmt.Sprintf("Score: %v", rgs.score)
		rl.DrawText(timeText, 500, 60, 40, rl.Black)

		rl.DrawText("Press enter to go back to main menu", 350, 200, 30, rl.Black)
	}

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
