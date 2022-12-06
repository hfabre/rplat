package main

import (
	rp "example.com/rplat/pkg/game"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	g := rp.NewGame()
	g.Run()
}
