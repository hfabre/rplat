package main

import (
	"runtime"
	rp "example.com/rplat/pkg/game"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	g := rp.NewGame()
	g.Run()
}