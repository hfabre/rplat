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
	player       *Player
	level        Map
	inputManager *InputManager
	stars        []Star
	score        int
	sceneManager *SceneManager
	gameEnded    bool
	helpOpen     bool
}

func NewTuorialGameScene(sm *SceneManager) *TuorialGameScene {
	tgs := &TuorialGameScene{}

	// Load level
	mc := NewMapConfiguration("./assets/map.json")
	tileset := NewTileset("./assets/tileset.png", mc.TileWidth, mc.TileHeight)
	im := NewInputManager()

	tgs.level = NewMap(mc, tileset)
	tgs.inputManager = &im
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

	for i := 0; i < 2; i++ {
		for !tgs.SpawnStar() {
		}
	}
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
			case "help":
				tgs.helpOpen = !tgs.helpOpen
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
			case "quit":
				tgs.gameEnded = true
			default:
				// Unknown event
			}
		}
	}
}

func (tgs TuorialGameScene) ShouldExit() bool {
	return false
}

func (tgs *TuorialGameScene) Update(deltaTime float32) {
	if tgs.gameEnded {
		tgs.sceneManager.SwapScene("main_menu")
	}

	if tgs.helpOpen {
		return
	}

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
			for !tgs.SpawnStar() {
			}
		}
	}

	for _, j := range starsToRemove {
		tgs.stars = RemoveIndexStar(tgs.stars, j)
	}
}

func (tgs TuorialGameScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	if tgs.helpOpen {
		tgs.DrawHelp()
	} else {
		tgs.DrawGame(factor)
	}
}

func (tgs TuorialGameScene) DrawHelp() {
	leftOffset := 150

	rl.DrawText("Help menu press H again to close.", 250, 50, 50, rl.Black)
	rl.DrawText("Press BACKSPACE to leave", 300, 110, 50, rl.Black)

	rl.DrawText("In random game mode you have 30 seconds to catch all the stars", int32(leftOffset), 200, 30, rl.Black)
	rl.DrawText("Your score depends on how much remaining time you still get.", int32(leftOffset), 230, 30, rl.Black)
	rl.DrawText("To achieve your mission you have access to multiple fast travel skills", int32(leftOffset), 270, 30, rl.Black)

	leftOffset = leftOffset - 130
	rl.DrawText("Teeworlds fan ? You can use a grappling hook using your mouse right click !", int32(leftOffset), 350, 30, rl.Black)
	rl.DrawText("Already played portal ? You can fire your portal gun using your mouse left click !", int32(leftOffset), 380, 30, rl.Black)
	rl.DrawText("And finally, you can dash in the direction you are going using left shift.", int32(leftOffset), 410, 30, rl.Black)
}

func (tgs TuorialGameScene) DrawGame(factor float64) {
	tgs.level.Draw()
	tgs.player.Draw(factor)
	for _, star := range tgs.stars {
		star.Draw()
	}

	scoreText := fmt.Sprintf("Score: %v", tgs.score)
	rl.DrawText(scoreText, 500, 60, 40, rl.Black)

	rl.DrawText("Press enter H to open help", 350, 200, 30, rl.Black)
}
