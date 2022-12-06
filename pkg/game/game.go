package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
)

const ScreenWidth = 1280
const ScreenHeight = 700
const FPS = 120

// Debug
const Debug = true

var Pause = false

type Game struct {
	currentTime  float64
	dt           float64
	time         float64
	accumulator  float64
	inputManager InputManager
	sm           SceneManager
}

func NewGame() Game {
	g := Game{}

	g.currentTime = rl.GetTime()
	g.dt = 0.01
	g.inputManager = NewInputManager()

	return g
}

func (g *Game) Run() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "rplat")
	defer rl.CloseWindow()

	scenes := make(map[string]Scene)
	g.sm = SceneManager{scenes: scenes, currentSceneName: "main_menu"}
	scenes["main_menu"] = NewMainMenuSceneWrapper(&g.sm)
	scenes["random_game"] = NewRandGameSceneWrapper(&g.sm)

	rl.SetTargetFPS(FPS)

	for !rl.WindowShouldClose() && !g.sm.ShouldExit() {
		g.Tick()
	}
}

// See https://gafferongames.com/post/fix_your_timestep/
func (g *Game) Tick() {
	newTime := rl.GetTime()
	frameTime := newTime - g.currentTime
	if frameTime > 0.25 {
		frameTime = 0.25
	}

	g.currentTime = newTime
	g.accumulator += frameTime

	for g.accumulator >= g.dt {
		g.sm.UpdateInputs()
		g.sm.HandleEvents()
		g.sm.Update(g.deltaTime(g.dt))

		g.time += g.dt
		g.accumulator -= g.dt
		g.sm.ClearInputs()
	}

	alpha := g.accumulator / g.dt
	g.sm.Draw(alpha)
}

func (g Game) deltaTime(dt float64) float32 {
	return float32(dt)
}
