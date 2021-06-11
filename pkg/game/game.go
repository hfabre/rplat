package game

import (
	"fmt"
	rl "github.com/chunqian/go-raylib/raylib"
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

type Game struct {
	player *Player
	level Map
	currentTime float64
	dt float64
	time float64
	accumulator float64
	inputManager InputManager
}

func NewGame() Game {
	player := Player{
		pos:      rl.Vector2{32, 32},
		lastPos:  rl.Vector2{20, 20},
		velocity: rl.Vector2{0, 0},
		lastVelocity: rl.Vector2{0, 0},
		size:     rl.Vector2{32, 64},
		canJump: true,
		color: rl.Red,
	}

	g := Game{player: &player, currentTime: rl.GetTime(), dt: 0.01, inputManager: NewInputManager()}

	return g
}

func (g *Game) Run() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "rplat")
	defer rl.CloseWindow()

	mc := NewMapConfiguration("./assets/map.json")
	tileset := NewTileset("./assets/tileset.png", mc.TileWidth, mc.TileHeight)
	g.level = NewMap(mc, tileset)

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
		g.inputManager.Update()
		g.HandleEvents()

		if !Pause {
			g.player.color = rl.Green
			g.player.lastPos = g.player.pos
			g.player.lastVelocity = g.player.velocity
			g.Update(g.deltaTime(g.dt))
		}

		g.time += g.dt
		g.accumulator -= g.dt
		g.inputManager.Clear()
	}

	alpha := g.accumulator / g.dt
	g.Draw(alpha)
}

func (g Game) deltaTime(dt float64) float32 {
	return float32(dt)
}

func (g *Game) HandleEvents() {
	for i := 0; i < len(g.inputManager.events); i++ {
		switch e := g.inputManager.events[i]; e {
		case "pause":
			Pause = !Pause
		case "move_right":
			g.player.velocity.X += PlayerSpeed
		case "move_left":
			g.player.velocity.X -= PlayerSpeed
		case "jump":
			if g.player.canJump {
				g.player.canJump = false
				g.player.velocity.Y -= PlayerJumpSpeed
			}
		case "hook":
			if !g.player.hookLaunched {
				g.player.hook = NewHook(*g.player)
				g.player.hookLaunched = true
			}
		case "stop_hook":
			g.player.hookLaunched = false
		default:
			// Unknown event
		}
	}
}

func (g *Game) Update(deltaTime float32) {
	g.player.Update(deltaTime)

	// Resolve collisions
	for i := 0; i < len(g.level.walls); i++ {
		if rl.CheckCollisionRecs(g.player.Rectangle(), g.level.walls[i]) {
			g.player.SolveCollision(g.level.walls[i])
		}

		if g.player.hookLaunched {
			if rl.CheckCollisionRecs(g.player.hook.Rectangle(), g.level.walls[i]) {
				g.player.hook.SolveCollision(g.level.walls[i])
			}
		}
	}
}

func (g Game) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	g.level.Draw()

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
