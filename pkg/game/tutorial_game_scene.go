package game

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/chunqian/go-raylib/raylib"
)

// Workaround to be able to call methods on a pointer on my interface
type TuorialGameSceneWrapper struct {
	tgs *TuorialGameScene
}

func NewTuorialGameSceneWrapper(sm *SceneManager) TuorialGameSceneWrapper {
	return TuorialGameSceneWrapper{tgs: NewTuorialGameScene(sm)}
}

// Implement Scene interface
func (tgsw TuorialGameSceneWrapper) Init() {
	tgsw.tgs.Init()
}

func (tgsw TuorialGameSceneWrapper) UpdateInputs() {
	tgsw.tgs.UpdateInputs()
}

func (tgsw TuorialGameSceneWrapper) ClearInputs() {
	tgsw.tgs.ClearInputs()
}

func (tgsw TuorialGameSceneWrapper) HandleEvents() {
	tgsw.tgs.HandleEvents()
}

func (tgsw TuorialGameSceneWrapper) Update(dt float32) {
	tgsw.tgs.Update(dt)
}

func (tgsw TuorialGameSceneWrapper) Draw(factor float64) {
	tgsw.tgs.Draw(factor)
}

func (tgsw TuorialGameSceneWrapper) End() {
	tgsw.tgs.End()
}

func (tgsw TuorialGameSceneWrapper) ShouldExit() bool {
	return tgsw.tgs.ShouldExit()
}

type TuorialGameScene struct {
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

func NewTuorialGameScene(sm *SceneManager) *TuorialGameScene {
	tgs := &TuorialGameScene{}

	// Load level
	mc := NewMapConfiguration("./assets/map.json")
	tileset := NewTileset("./assets/tileset.png", mc.TileWidth, mc.TileHeight)
	im := NewInputManager()

	tgs.level = NewMap(mc, tileset)
	tgs.inputManager = &im
	tgs.durationSeconds = 30
	tgs.sceneManager = sm

	return tgs
}

func (tgs *TuorialGameScene) Init() {
	player := Player{
		pos:          rl.Vector2{X: 32, Y: 32},
		lastPos:      rl.Vector2{X: 20, Y: 20},
		velocity:     rl.Vector2{X: 0, Y: 0},
		lastVelocity: rl.Vector2{X: 0, Y: 0},
		size:         rl.Vector2{X: 32, Y: 64},
		canJump:      true,
		color:        rl.Red,
	}

	tgs.player = &player
	tgs.gameEnded = false
	tgs.ticker = time.NewTicker(1 * time.Second)
	tgs.elapsedSeconds = 0

	for i := 0; i < 20; i++ {
		for !tgs.SpawnStar() {
		}
	}

	go tgs.StartsCounter()
}

func (tgs *TuorialGameScene) End() {
	tgs.stars = nil
}

func (tgs *TuorialGameScene) SpawnStar() bool {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	x := r.Intn(ScreenWidth)
	y := r.Intn(ScreenHeight)

	star := Star{rl.Vector2{X: float32(x), Y: float32(y)}}

	for _, wall := range tgs.level.walls {
		if isColliding(star.Rectangle(), wall) {
			return false
		}
	}

	for _, exStar := range tgs.stars {
		if isColliding(star.Rectangle(), exStar.Rectangle()) {
			return false
		}
	}

	tgs.stars = append(tgs.stars, star)
	return true
}

func (tgs *TuorialGameScene) UpdateInputs() {
	tgs.inputManager.Update()
}

func (tgs *TuorialGameScene) ClearInputs() {
	tgs.inputManager.Clear()
}

func (tgs *TuorialGameScene) HandleEvents() {
	for i := 0; i < len(tgs.inputManager.events); i++ {
		if tgs.gameEnded && tgs.inputManager.events[i] == "validate" {
			tgs.sceneManager.SwapScene("main_menu")
		} else {
			switch e := tgs.inputManager.events[i]; e {
			case "pause":
				Pause = !Pause
			case "move_right":
				tgs.player.MoveRight()
			case "move_left":
				tgs.player.MoveLeft()
			case "jump":
				tgs.player.Jump()
			case "hook":
				tgs.player.Hook()
			case "stop_hook":
				tgs.player.StopHook()
			case "dash":
				tgs.player.Dash()
			case "portal":
				tgs.player.FirePortal(tgs.level.walls)
			default:
				// Unknown event
			}
		}
	}
}

func (tgs *TuorialGameScene) StartsCounter() {
	for range tgs.ticker.C {
		if !Pause {
			tgs.elapsedSeconds++

			if tgs.elapsedSeconds >= tgs.durationSeconds {
				tgs.EndGame(false)
				return
			}
		}
	}
}

func (tgs *TuorialGameScene) EndGame(complete bool) {
	if complete {
		multiplier := tgs.durationSeconds - tgs.elapsedSeconds
		tgs.score *= multiplier
	}

	tgs.gameEnded = true
	tgs.ticker.Stop()
}

func (tgs TuorialGameScene) ShouldExit() bool {
	return false
}

func (tgs *TuorialGameScene) Update(deltaTime float32) {
	if tgs.gameEnded {
		return
	}

	if !Pause {
		tgs.player.color = rl.Green
		tgs.player.lastPos = tgs.player.pos
		tgs.player.lastVelocity = tgs.player.velocity

		tgs.player.Update(deltaTime)
		tgs.player.checkAndHandleCollisions(tgs.level.walls)

		starsToRemove := []int{}
		for i, star := range tgs.stars {
			if isColliding(tgs.player.Rectangle(), star.Rectangle()) {
				tgs.score += 10
				starsToRemove = append(starsToRemove, i)
			}
		}

		for _, j := range starsToRemove {
			tgs.stars = RemoveIndexStar(tgs.stars, j)
		}

		if len(tgs.stars) == 0 {
			tgs.EndGame(true)
		}
	}
}

func (tgs TuorialGameScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	tgs.level.Draw()
	tgs.player.Draw(factor)
	for _, star := range tgs.stars {
		star.Draw()
	}

	timeText := fmt.Sprintf("Elapsed time: %v", tgs.elapsedSeconds)
	rl.DrawText(timeText, 500, 20, 40, rl.Black)
	scoreText := fmt.Sprintf("Score: %v", tgs.score)
	rl.DrawText(scoreText, 500, 60, 40, rl.Black)

	rl.DrawText("Press enter H to open help", 350, 200, 30, rl.Black)
}
