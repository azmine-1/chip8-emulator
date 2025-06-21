package main

import (
	"fmt"
	"os"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)
const Font_data = []byte{0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80} // F

var display_grid [64][32]bool
type Stack struct {
	data []uint16
}


var KeyPad = [4][4]string{
	{"1", "2", "3", "C"},
	{"4", "5", "6", "D"},
	{"7", "8", "9", "E"},
	{"A", "0", "B", "F"},
}


var KeyMap = map[string]byte{
	"1": 0x1, "2": 0x2, "3": 0x3, "C": 0xC,
	"4": 0x4, "5": 0x5, "6": 0x6, "D": 0xD,
	"7": 0x7, "8": 0x8, "9": 0x9, "E": 0xE,
	"A": 0xA, "0": 0x0, "B": 0xB, "F": 0xF,
}

func(s *Stack) push(value uint16){
	if (len(s.data)) > 16{
		return nil
	}
	else{
		s.data = append(s.data, value)
	}
}
func(s *Stack) pop(value uint16, b bool){
	if(len(s.data)  == 0){
		return 0, false
	}
	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
}
var delayTimer byte
var soundTimer byte

func start_timers(){
	ticker = time.NewTicker(time.second /60)
	go func(){
		for range ticker.C{
			if(delayTimer > 0){
				delayTimer--
			}
			if(soundTimer > 0){
				soundTimer--;
			}
		}
	}()
}

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
	return 64, 32
}


