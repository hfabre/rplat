package game

import (
	rl "github.com/chunqian/go-raylib/raylib"
)

type MenuInputManager struct {
	inputMap map[string]int32
	events   []string
}

func NewMenuInputManager() MenuInputManager {
	im := MenuInputManager{}
	m := make(map[string]int32)

	m["move_left"] = int32(rl.KEY_A)
	m["move_left"] = int32(rl.KEY_LEFT)
	m["move_right"] = int32(rl.KEY_D)
	m["move_right"] = int32(rl.KEY_RIGHT)
	m["move_up"] = int32(rl.KEY_W)
	m["move_up"] = int32(rl.KEY_UP)
	m["move_down"] = int32(rl.KEY_S)
	m["move_down"] = int32(rl.KEY_DOWN)
	m["validate"] = int32(rl.KEY_ENTER)

	im.inputMap = m
	return im
}

func (im *MenuInputManager) Update() {
	if rl.IsKeyPressed(im.inputMap["move_up"]) {
		im.events = append(im.events, "move_up")
	}

	if rl.IsKeyPressed(im.inputMap["move_down"]) {
		im.events = append(im.events, "move_down")
	}

	if rl.IsKeyPressed(im.inputMap["move_left"]) {
		im.events = append(im.events, "move_left")
	}

	if rl.IsKeyPressed(im.inputMap["move_right"]) {
		im.events = append(im.events, "move_right")
	}

	if rl.IsKeyPressed(im.inputMap["validate"]) {
		im.events = append(im.events, "validate")
	}
}

func (im *MenuInputManager) Clear() {
	im.events = nil
}
