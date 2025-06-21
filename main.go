package main

import (
	"time"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	chip := chip8.New()
	chip.initialize()
	chip.loadROM

	for {
		chip.EmulateCycle()
		chip.HandleInput()
		chip.render()
		time.Sleep(time.Miillisecond * 2)
	}
}
