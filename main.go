package main

import (
	"chip8"
	"time"
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
