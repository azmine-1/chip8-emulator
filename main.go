package main

import (
	"fmt"
	"os"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Println("Starting CHIP-8 Emulator...")
	
	// Initialize memory
	memory := &Memory{
		memory: make([]byte, 4096),
		PC:     0x200, // CHIP-8 programs start at 0x200
	}
	
	// Load font data into memory
	for i, fontByte := range Font_data {
		memory.memory[i] = fontByte
	}
	
	// Initialize stack
	stack := &Stack{
		data: make([]uint16, 0, 16),
	}
	
	// Start timers
	start_timers()
	
	// Set up the game window
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	
	// Create game instance
	game := &Game{
		memory: memory,
		stack:  stack,
	}
	
	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}
}
