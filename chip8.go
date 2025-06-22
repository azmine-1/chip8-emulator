package main

import (
	"fmt"
	"os"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)
var Font_data = []byte{0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
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
type Memory struct{
	memory []byte 
	PC uint16
}
func fetch(m *Memory) uint16 {
	var cur_instruction [2]byte
	for i := 0; i < 2; i++ {
		cur_instruction[i] = m.memory[m.PC + uint16(i)]
	}
	m.PC += 2
	return uint16(cur_instruction[0])<<8 | uint16(cur_instruction[1])
}
func decode(opcode uint16) {
    if opcode == 0x00E0 {
        fmt.Println("CLS")
    } else if opcode == 0x00EE {
        fmt.Println("RET")
    } else if opcode&0xF000 == 0x1000 {
        addr := opcode & 0x0FFF
        fmt.Printf("JP 0x%03X\n", addr)
    } else if opcode&0xF000 == 0x2000 {
        addr := opcode & 0x0FFF
        fmt.Printf("CALL 0x%03X\n", addr)
    } else if opcode&0xF000 == 0x3000 {
        x := (opcode & 0x0F00) >> 8
        kk := opcode & 0x00FF
        fmt.Printf("SE V%X, 0x%02X\n", x, kk)
    } else if opcode&0xF000 == 0x6000 {
        x := (opcode & 0x0F00) >> 8
        kk := opcode & 0x00FF
        fmt.Printf("LD V%X, 0x%02X\n", x, kk)
    } else if opcode&0xF000 == 0x7000 {
        x := (opcode & 0x0F00) >> 8
        kk := opcode & 0x00FF
        fmt.Printf("ADD V%X, 0x%02X\n", x, kk)
    } else if opcode&0xF000 == 0xA000 {
        addr := opcode & 0x0FFF
        fmt.Printf("LD I, 0x%03X\n", addr)
    } else if opcode&0xF000 == 0xD000 {
        x := (opcode & 0x0F00) >> 8
        y := (opcode & 0x00F0) >> 4
        n := opcode & 0x000F
        fmt.Printf("DRW V%X, V%X, %X\n", x, y, n)
    } else {
        fmt.Printf("Unknown opcode: 0x%04X\n", opcode)
    }
}
func execute(opcode uint16){
	decode(opcode);
}

var KeyPad = [4][4]string{
	{"1", "2", "3", "C"},
	{"4", "5", "6", "D"},
	{"7", "8", "9", "E"},
	{"A", "0", "B", "F"},
}

var KeyPad_KeyBoard = [4][4]string{
	{"1", "2", "3", "4"},
	{"Q" ,"W", "E", "R"},
	{"A", "S", "D", "F"},
	{"Z", "X", "C", "V"},
}


var KeyMap = map[string]byte{
	"1": 0x1, "2": 0x2, "3": 0x3, "C": 0xC,
	"4": 0x4, "5": 0x5, "6": 0x6, "D": 0xD,
	"7": 0x7, "8": 0x8, "9": 0x9, "E": 0xE,
	"A": 0xA, "0": 0x0, "B": 0xB, "F": 0xF,
}
var KeyMap_KeyBoard = map[string]byte{
	"1": 0x1, "2": 0x2, "3": 0x3, "4": 0xC,
	"Q": 0x4, "W": 0x5, "E": 0x6, "R": 0xD,
	"A": 0x7, "S": 0x8, "D": 0x9, "F": 0xE,
	"Z": 0xA, "X": 0x0, "C": 0xB, "V": 0xF,
}

func(s *Stack) push(value uint16) error {
	if len(s.data) >= 16 {
		return fmt.Errorf("stack overflow")
	}
	s.data = append(s.data, value)
	return nil
}

func(s *Stack) pop() (uint16, error) {
	if len(s.data) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}
	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return value, nil
}

var delayTimer byte
var soundTimer byte
var ticker *time.Ticker

func start_timers() {
	ticker = time.NewTicker(time.Second / 60)
	go func() {
		for range ticker.C {
			if delayTimer > 0 {
				delayTimer--
			}
			if soundTimer > 0 {
				soundTimer--
			}
		}
	}()
}

type Game struct {
	memory *Memory
	stack  *Stack
}

func(g *Game) Update() error {
	opcode := fetch(g.memory)
	execute(opcode)
	return nil
}

func(g *Game) Draw(screen *ebiten.Image){
	screen.Fill(color.RGBA{0, 0, 255,255})
	ebitenutil.DebugPrint(screen, "Chip8-emulator")
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int){
	return 64, 32
}


