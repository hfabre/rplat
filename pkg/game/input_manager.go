package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
	"time"
)

var TimeBetweenPauses = 100 * time.Millisecond
var lastPauseTime = time.Now()

type InputManager struct {
	inputMap map[string]int32
	events []string
}

func NewInputManager() InputManager {
	im := InputManager{}
	m := make(map[string]int32)

	m["jump"] = int32(rl.KEY_SPACE)
	m["move_left"] = int32(rl.KEY_A)
	m["move_right"] = int32(rl.KEY_D)
	m["hook"] = int32(rl.KEY_ENTER)
	m["mouse_hook"] = int32(rl.MOUSE_RIGHT_BUTTON)

	im.inputMap = m
	return im
}

func (im *InputManager) Update() {
	kbHook := false

	if Debug && time.Since(lastPauseTime) > TimeBetweenPauses  && rl.IsKeyDown(int32(rl.KEY_P)) {
		lastPauseTime = time.Now()
		im.events = append(im.events, "pause")
	}

	if !Pause {
		if rl.IsKeyDown(im.inputMap["move_left"]) {
			im.events = append(im.events, "move_left")
		}

		if rl.IsKeyDown(im.inputMap["move_right"]) {
			im.events = append(im.events, "move_right")
		}

		if rl.IsKeyDown(im.inputMap["jump"]) {
			im.events = append(im.events, "jump")
		}

		if rl.IsKeyDown(im.inputMap["hook"]) {
			im.events = append(im.events, "hook")
			kbHook = true
		}

		if rl.IsKeyUp(im.inputMap["hook"]) && kbHook {
			im.events = append(im.events, "stop_hook")
		}

		if rl.IsMouseButtonDown(im.inputMap["mouse_hook"]) {
			im.events = append(im.events, "hook")
		}

		if rl.IsMouseButtonUp(im.inputMap["mouse_hook"]) && !kbHook {
			im.events = append(im.events, "stop_hook")
		}
	}
}

func (im *InputManager) Clear() {
	im.events = nil
}
