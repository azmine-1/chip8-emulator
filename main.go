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
		memory.memory[i] = fontByte
	}
	
	
	stack := &Stack{
		data: make([]uint16, 0, 16),
	}
	
	
	start_timers()
	
	
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	
	
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
