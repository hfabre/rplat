package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
)

// Workaround to be able to call methods on a pointer on my interface
type MainMenuSceneWrapper struct {
	mms *MainMenuScene
}

func NewMainMenuSceneWrapper() MainMenuSceneWrapper {
	return MainMenuSceneWrapper{mms: NewMainMenuScene()}
}

// Implement Scene interface
func (mmsw MainMenuSceneWrapper) Init() {

}

func (mmsw MainMenuSceneWrapper) UpdateInputs() {
	mmsw.mms.UpdateInputs()
}

func (mmsw MainMenuSceneWrapper) ClearInputs() {
	mmsw.mms.ClearInputs()
}

func (mmsw MainMenuSceneWrapper) HandleEvents() {
	mmsw.mms.HandleEvents()
}

func (mmsw MainMenuSceneWrapper) Update(dt float32) {
	mmsw.mms.Update(dt)
}

func (mmsw MainMenuSceneWrapper) Draw(factor float64) {
	mmsw.mms.Draw(factor)
}

func (mmsw MainMenuSceneWrapper) End() {
	mmsw.mms.End()
}

type MainMenuScene struct {
	inputManager *MenuInputManager
	selectedItem int
	items        []string
}

func NewMainMenuScene() *MainMenuScene {
	mms := &MainMenuScene{}

	im := NewMenuInputManager()
	mms.inputManager = &im

	mms.items = append(mms.items, "Tutorial")
	mms.items = append(mms.items, "Random game")
	mms.items = append(mms.items, "Exit")

	return mms
}

func (mms *MainMenuScene) UpdateInputs() {
	mms.inputManager.Update()
}

func (mms *MainMenuScene) ClearInputs() {
	mms.inputManager.Clear()
}

func (mms *MainMenuScene) End() {

}

func (mms *MainMenuScene) HandleEvents() {
	for i := 0; i < len(mms.inputManager.events); i++ {
		switch e := mms.inputManager.events[i]; e {
		case "move_right":
		case "move_left":
		case "move_up":
		case "move_down":
		case "validate":
		default:
			// Unknown event
		}
	}
}

func (mms *MainMenuScene) Update(deltaTime float32) {

}

func (mms MainMenuScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()
	var color rl.Color
	var x int32
	var y int32

	for i, item := range mms.items {
		color = rl.Black
		if i == mms.selectedItem {
			color = rl.Green
		}

		x = int32(ScreenWidth/2 - (20 + i + 1))
		y = int32(ScreenHeight/2 - (20 + i + 1))
		rl.DrawText(item, x, y, 20, color)
	}

	rl.ClearBackground(rl.RayWhite)
}
