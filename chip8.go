package main

import (
	"fmt"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"os"
	"image"
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
type Memory struct {
	memory []byte
	PC     uint16
	V      [16]byte 
	I      uint16   
}
func fetch(m *Memory) uint16 {
	var cur_instruction [2]byte
	for i := 0; i < 2; i++ {
		cur_instruction[i] = m.memory[m.PC + uint16(i)]
	}
	m.PC += 2
	return uint16(cur_instruction[0])<<8 | uint16(cur_instruction[1])
}
func decode(opcode uint16, m *Memory, s *Stack) {
	switch {
	case opcode == 0x00E0:
		for x := range display_grid {
			for y := range display_grid[x] {
				display_grid[x][y] = false
			}
		}
	case opcode == 0x00EE:
		retAddr, err := s.pop()
		if err == nil {
			m.PC = retAddr
		}
	case opcode&0xF000 == 0x1000:
		addr := opcode & 0x0FFF
		m.PC = addr

	case opcode&0xF000 == 0x2000:
		addr := opcode & 0x0FFF
		s.push(m.PC)
		m.PC = addr

	case opcode&0xF000 == 0x3000:
		x := (opcode & 0x0F00) >> 8
		kk := byte(opcode & 0x00FF)
		if m.V[x] == kk {
			m.PC += 2
		}

	case opcode&0xF000 == 0x4000:
		x := (opcode & 0x0F00) >> 8
		kk := byte(opcode & 0x00FF)
		if m.V[x] != kk {
			m.PC += 2
		}

	case opcode&0xF00F == 0x5000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if m.V[x] == m.V[y] {
			m.PC += 2
		}

	case opcode&0xF000 == 0x6000:
		x := (opcode & 0x0F00) >> 8
		kk := byte(opcode & 0x00FF)
		m.V[x] = kk

	case opcode&0xF000 == 0x7000:
		x := (opcode & 0x0F00) >> 8
		kk := byte(opcode & 0x00FF)
		m.V[x] += kk

	case opcode&0xF000 == 0x8000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		switch opcode & 0x000F {
		case 0x0:
			m.V[x] = m.V[y]
		case 0x1:
			m.V[x] |= m.V[y]
		case 0x2:
			m.V[x] &= m.V[y]
		case 0x3:
			m.V[x] ^= m.V[y]
		case 0x4:
			sum := uint16(m.V[x]) + uint16(m.V[y])
			m.V[0xF] = 0
			if sum > 255 {
				m.V[0xF] = 1
			}
			m.V[x] = byte(sum)
		case 0x5:
			m.V[0xF] = 0
			if m.V[x] > m.V[y] {
				m.V[0xF] = 1
			}
			m.V[x] -= m.V[y]
		case 0x6:
			m.V[0xF] = m.V[x] & 0x1
			m.V[x] >>= 1
		case 0x7:
			m.V[0xF] = 0
			if m.V[y] > m.V[x] {
				m.V[0xF] = 1
			}
			m.V[x] = m.V[y] - m.V[x]
		case 0xE:
			m.V[0xF] = (m.V[x] & 0x80) >> 7
			m.V[x] <<= 1
		}

	case opcode&0xF000 == 0xA000:
		m.I = opcode & 0x0FFF

	case opcode&0xF000 == 0xD000:
		x := uint16(m.V[(opcode&0x0F00)>>8] % 64)
		y := uint16(m.V[(opcode&0x00F0)>>4] % 64)
		height := opcode & 0x000F
		m.V[0xF] = 0
		for row := uint16(0); row < height; row ++{
			spriteByte := m.memory[m.I+row]
			for col := uint16(0); col < 8; col++{
				if(spriteByte & (0x80 >> col)) != 0{
					px := (x + col) % 64
					py := (y + row ) % 32

					if display_grid[int(px)][int(py)] {
						m.V[0xF] = 1
					}
					display_grid[int(px)][int(py)] = !display_grid[int(px)][int(py)]
				}

			}
		}

	default:
		fmt.Printf("Unknown opcode: 0x%04X\n", opcode)
	}
}
func execute(opcode uint16, m *Memory, s *Stack) {
	decode(opcode, m, s)
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

// Load test program into memory
func loadTestProgram(m *Memory){
	// Instructions
	testInstructions := []byte{
		0x00, 0xE0,       // CLS
		0x60, 0x00,       // LD V0, 0x00
		0x61, 0x00,       // LD V1, 0x00
		0xA2, 0x2A,       // LD I, 0x22A
		0xD0, 0x1F,       // DRW V0, V1, 0xF
		0x70, 0x08,       // ADD V0, 0x08
		0xA2, 0x39,       // LD I, 0x239
		0xD0, 0x1F,       // DRW V0, V1, 0xF
		0x70, 0x08,       // ADD V0, 0x08
		0xA2, 0x48,       // LD I, 0x248
		0xD0, 0x1F,       // DRW V0, V1, 0xF
		0x70, 0x08,       // ADD V0, 0x08
		0xA2, 0x57,       // LD I, 0x257
		0xD0, 0x1F,       // DRW V0, V1, 0xF
		0x12, 0x20,       // JP 0x220
	}

	// Load program into memory starting at 0x200
	for i, instruction := range testInstructions {
		m.memory[0x200+i] = instruction
	}

	// Sprite data for letters "I", "B", "M" and a stylized block
	ibmLogo := []byte{
		// Sprite 1 - "I"
		0xFF, 0x18, 0x18, 0x18, 0x18,
		0x18, 0x18, 0x18, 0x18, 0x18,
		0x18, 0x18, 0x18, 0x18, 0xFF,

		// Sprite 2 - "B"
		0xFE, 0x92, 0x92, 0x92, 0xFE,
		0x92, 0x92, 0x92, 0x92, 0xFE,
		0x92, 0x92, 0x92, 0x92, 0xFE,

		// Sprite 3 - "M"
		0x81, 0xC3, 0xA5, 0xA5, 0x99,
		0x99, 0x81, 0x81, 0x81, 0x81,
		0x81, 0x81, 0x81, 0x81, 0x81,

		// Sprite 4 - Solid bar (or spacing block)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	// Load into memory starting at 0x22A
	for i, spriteByte := range ibmLogo {
		m.memory[0x22A+i] = spriteByte
	}
}


func loadROM(m *Memory, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %v", err)
	}
	
	// Load ROM data starting at 0x20
	for i, byte := range data {
		if 0x200+i >= 4096 {
			break 
		}
		m.memory[0x200+i] = byte
	}
	
	return nil
}

type Game struct {
	memory *Memory
	stack  *Stack
}

func(g *Game) Update() error {
	opcode := fetch(g.memory)
	execute(opcode, g.memory, g.stack)
	return nil
}

func(g *Game) Draw(screen *ebiten.Image){

	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			if display_grid[x][y] {
				rect := image.Rect(x*10, y*10, (x+1)*10, (y+1)*10)
				ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), 
					float64(rect.Dx()), float64(rect.Dy()), color.RGBA{255, 255, 255, 255})
			}
		}
	}
	
	
	ebitenutil.DebugPrint(screen, fmt.Sprintf("PC: 0x%03X I: 0x%03X V0: 0x%02X V1: 0x%02X", 
		g.memory.PC, g.memory.I, g.memory.V[0], g.memory.V[1]))
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int){
	return 640, 320
}


