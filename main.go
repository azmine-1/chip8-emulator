package main

import (
	"fmt"
	"os"
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
	
	// Load program based on command line arguments
	if len(os.Args) > 1 {
		// Load ROM file if provided
		romFile := os.Args[1]
		fmt.Printf("Loading ROM: %s\n", romFile)
		if err := loadROM(memory, romFile); err != nil {
			fmt.Printf("Error loading ROM: %v\n", err)
			fmt.Println("Loading test program instead...")
			loadTestProgram(memory)
		}
	} else {
		// Load test program by default
		fmt.Println("No ROM specified, loading test program...")
		loadTestProgram(memory)
	}
	
	start_timers()
	
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	
	game := &Game{
		memory: memory,
		stack:  stack,
	}
	
	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}
}
