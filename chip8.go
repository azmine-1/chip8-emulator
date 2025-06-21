package main

import (
	"fmt"
	"os"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

func main(){
	fmt.Println("works")
	os.Getenv("HOME")
	ebiten.SetWindowSize(640,480)
	ebiten.SetWindowTitle("Chip8-emulator")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

type Game struct{}

func(g *Game) Update() error {
	return nil;
}

func(g *Game) Draw(screen *ebiten.Image){
	screen.Fill(color.RGBA{0, 0, 255,255})
	ebitenutil.DebugPrint(screen, "Chip8-emulator")
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int){
	return 320, 240
}


