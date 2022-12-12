package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
)

// Workaround to be able to call methods on a pointer on my interface
type MainMenuSceneWrapper struct {
	mms *MainMenuScene
}

func NewMainMenuSceneWrapper(sm *SceneManager) MainMenuSceneWrapper {
	return MainMenuSceneWrapper{mms: NewMainMenuScene(sm)}
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

func (mmsw MainMenuSceneWrapper) ShouldExit() bool {
	return mmsw.mms.ShouldExit()
}

type MainMenuScene struct {
	inputManager *MenuInputManager
	selectedItem int
	items        []string
	sceneManager *SceneManager
	exit         bool
}

func NewMainMenuScene(sm *SceneManager) *MainMenuScene {
	mms := &MainMenuScene{}

	im := NewMenuInputManager()
	mms.inputManager = &im
	mms.sceneManager = sm
	mms.exit = false
	mms.selectedItem = 1

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
		case "move_up":
			if mms.selectedItem <= 0 {
				mms.selectedItem = len(mms.items) - 1
			} else {
				mms.selectedItem -= 1
			}
		case "move_down":
			if mms.selectedItem >= len(mms.items)-1 {
				mms.selectedItem = 0
			} else {
				mms.selectedItem += 1
			}
		case "validate":
			switch mms.selectedItem {
			case 0:
				mms.sceneManager.SwapScene("tutorial_game")
			case 1:
				mms.sceneManager.SwapScene("random_game")
			case 2:
				mms.exit = true
			}
		default:
			// Unknown menu item
		}
	}
}

func (mms *MainMenuScene) Update(deltaTime float32) {

}

func (mms MainMenuScene) ShouldExit() bool {
	return mms.exit
}

func (mms MainMenuScene) Draw(factor float64) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.DrawText(mms.items[0], 500, 100, 20, mms.ColorFromItem(0))
	rl.DrawText(mms.items[1], 500, 130, 20, mms.ColorFromItem(1))
	rl.DrawText(mms.items[2], 500, 160, 20, mms.ColorFromItem(2))

	rl.ClearBackground(rl.RayWhite)
}

func (mms MainMenuScene) ColorFromItem(item_index int) rl.Color {
	if item_index == mms.selectedItem {
		return rl.Green
	} else {
		return rl.Black
	}
}
