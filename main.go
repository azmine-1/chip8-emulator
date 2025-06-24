package main

import (
	"fmt"
	"os"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Println("Starting CHIP-8 Emulator...")
	
	memory := &Memory{
		memory: make([]byte, 4096),
		PC:     0x200, 
	}
	
	for i, fontByte := range Font_data {
		memory.memory[0x50 + i] = fontByte
	}
	
	stack := &Stack{
		data: make([]uint16, 0, 16),
	}
	if len(os.Args) > 1 {
		romFile := os.Args[1]
		fmt.Printf("Loading ROM: %s\n", romFile)
		if err := loadROM(memory, romFile); err != nil {
			fmt.Printf("Error loading ROM: %v\n", err)
			fmt.Println("Loading test program instead...")
			loadTestProgram(memory)
		}
	} else {
		fmt.Println("No ROM specified, loading test program...")
		loadTestProgram(memory)
	}
	
	start_timers()
	
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	
	game := &Game{
		memory: memory,
		stack:  stack,
		lastUpdate: time.Now(),
		instructionsPerFrame: 0,
	}
	
	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}
}
