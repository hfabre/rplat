package game

import (
	"fmt"
	rl "github.com/chunqian/go-raylib/raylib"
	"time"
)


const ScreenWidth = 1280
const ScreenHeight = 720
const FPS = 500

// In pixel per second
const PlayerSpeed = 100
const PlayerJumpSpeed = 500
const HookSpeed = 1800
const HookVerticalForce = 30
const HookHorizontalForce = 60

const Friction = 0.80
const Gravity = 10

// Debug
const Debug = true

var Pause = false
var TimeBetweenPauses = 100 * time.Millisecond
var lastPauseTime = time.Now()

var Slow = false
var TimeBetweenSlow = 100 * time.Millisecond
var lastSlowTime = time.Now()

type Game struct {
	player *Player
	walls []rl.Rectangle
	currentTime float64
	dt float64
	time float64
	accumulator float64
}

func NewGame() Game {
	player := Player{
		pos:      rl.Vector2{20, 20},
		lastPos:  rl.Vector2{20, 20},
		velocity: rl.Vector2{0, 0},
		lastVelocity: rl.Vector2{0, 0},
		size:     rl.Vector2{32, 64},
		canJump: true,
		color: rl.Green,
	}

	walls := make([]rl.Rectangle, 3)
	walls[0] = rl.Rectangle{0, 688, 1280, 32}
	walls[1] = rl.Rectangle{500, 530, 280, 32}
	walls[2] = rl.Rectangle{0, 200, 500, 32}

	g := Game{player: &player, walls: walls, currentTime: rl.GetTime(), dt: 0.01}

	return g
}

func (g Game) Run() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "rplat")
	defer rl.CloseWindow()

	rl.SetTargetFPS(FPS)

	for !rl.WindowShouldClose() {
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
		g.HandleInput(g.deltaTime(g.dt))

		if !Pause {
			g.player.color = rl.Green
			g.player.lastPos = g.player.pos
			g.player.lastVelocity = g.player.velocity
			g.Update(g.deltaTime(g.dt))
		}

		g.time += g.dt
		g.accumulator -= g.dt
	}

	alpha := g.accumulator / g.dt
	g.Draw(alpha)
}

func (g Game) deltaTime(dt float64) float32 {
	return float32(dt)
}

func (g *Game) HandleInput(deltaTime float32) {
	if Debug && time.Since(lastPauseTime) > TimeBetweenPauses  && rl.IsKeyDown(int32(rl.KEY_P)) {
		lastPauseTime = time.Now()
		Pause = !Pause
	}

	if Debug && time.Since(lastSlowTime) > TimeBetweenSlow  && rl.IsKeyDown(int32(rl.KEY_M)) {
		lastSlowTime = time.Now()
		Slow = !Slow
	}

	if !Pause {
		if rl.IsKeyDown(int32(rl.KEY_A)) {
			g.player.velocity.X -= PlayerSpeed
		}

		if rl.IsKeyDown(int32(rl.KEY_D)) {
			g.player.velocity.X += PlayerSpeed
		}

		if g.player.canJump && rl.IsKeyDown(int32(rl.KEY_SPACE)) {
			g.player.canJump = false
			g.player.velocity.Y -= PlayerJumpSpeed
		}

		if rl.IsKeyDown(int32(rl.KEY_ENTER)) {
			if !g.player.hookLaunched {
				g.player.hook = NewHook(*g.player)
				g.player.hookLaunched = true
			}
		}

		if rl.IsKeyUp(int32(rl.KEY_ENTER)) {
			g.player.hookLaunched = false
		}
	}
}

func (g *Game) Update(deltaTime float32) {
	g.player.Update(deltaTime)

	// Resolve collisions
	for i := 0; i < len(g.walls); i++ {
		if rl.CheckCollisionRecs(g.player.Rectangle(), g.walls[i]) {
			g.player.SolveCollision(g.walls[i])
		}

		if g.player.hookLaunched {
			if rl.CheckCollisionRecs(g.player.hook.Rectangle(), g.walls[i]) {
				g.player.hook.SolveCollision(g.walls[i])
			}
		}
	}
}

func (g Game) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	for i := 0; i < len(g.walls); i++ {
		rl.DrawRectangleRec(g.walls[i], rl.Gray)
	}

	currentStateLerp := LerpVec2(g.player.pos, factor)
	lastStateLerp := LerpVec2(g.player.pos, 1 - factor)
	rl.DrawRectangleV(rl.Vector2{currentStateLerp.X + lastStateLerp.X, currentStateLerp.Y + lastStateLerp.Y}, g.player.size, g.player.color)

	if g.player.hookLaunched {
		currentStateLerp = LerpVec2(g.player.hook.pos, factor)
		lastStateLerp = LerpVec2(g.player.hook.pos, 1 - factor)
		rl.DrawRectangleV(rl.Vector2{currentStateLerp.X + lastStateLerp.X, currentStateLerp.Y + lastStateLerp.Y}, g.player.hook.size, g.player.hook.color)

		rl.DrawLineEx(g.player.pos, g.player.hook.pos, 5, rl.Black)
	}

	if Debug {
		if Pause {
			rl.DrawRectangleV(rl.Vector2{g.player.collision.X, g.player.collision.Y}, rl.Vector2{10, 10}, rl.Blue)
		}

		posText := fmt.Sprintf("Position: %v - %v", g.player.pos.X, g.player.pos.Y)
		lastPosText := fmt.Sprintf("Last position: %v - %v", g.player.lastPos.X, g.player.lastPos.Y)
		velText := fmt.Sprintf("Velocity: %v - %v", g.player.velocity.X, g.player.velocity.Y)
		lastVelText := fmt.Sprintf("Last Velocity: %v - %v", g.player.lastVelocity.X, g.player.lastVelocity.Y)
		colText := fmt.Sprintf("Last collision: %v", g.player.collision)

		rl.DrawFPS(10, 10)
		rl.DrawText(posText, 10, 50, 20, rl.Black)
		rl.DrawText(lastPosText, 10, 70, 20, rl.Black)
		rl.DrawText(velText, 10, 90, 20, rl.Black)
		rl.DrawText(lastVelText, 10, 110, 20, rl.Black)
		rl.DrawText(colText, 10, 130, 20, rl.Black)
	}
}
