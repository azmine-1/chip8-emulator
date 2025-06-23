package main

import (
	"fmt"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"os"
	"image"
	"math/rand"
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
		y := uint16(m.V[(opcode&0x00F0)>>4] % 32)
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
	case opcode&0xF00F == 0x9000:
		x := (opcode&0x0F00) >> 8
		y := (opcode&0x00F0) >> 4
		if m.V[x] != m.V[y]{
			m.PC += 2
		}
	case opcode&0xF000 == 0xB000:
		addr := opcode&0x0FFF
		m.PC = addr + uint16(m.V[0])
	
	case opcode&0xF000 == 0xC000:
		x := (opcode&0x0F00) >> 8
		kk := byte(opcode&0x00FF)
		m.V[x] = byte(rand.Intn(256)) & kk
	
	case opcode&0xF000 == 0xE09E:
		x := (opcode&0xF000) >> 8
		key := m.V[x]
		if isKeyPressed(key){
			m.PC += 2
		}
	case opcode&0xF0FF == 0xE0A1:
		x := (opcode & 0x0F00) >> 8
		key := m.V[x]
    	if !isKeyPressed(key) {
        	m.PC += 2
    	}
	case opcode&0xF0FF == 0xF007:
		x := (opcode & 0x0F00) >> 8
		m.V[x] = delayTimer
	case opcode&0xF0FF == 0xF00A:
		x := (opcode&0x0F00) >> 8
		key := getKeyPressed()
		if key != 0xFF {
			m.V[x] = key
		} else {
			m.PC -= 2
		}
	case opcode&0xF0FF == 0xF015:
		x := (opcode & 0x0F00) >> 8
		delayTimer = m.V[x]
	case opcode&0xF0FF == 0xF018:
		x := (opcode&0xF00) >> 8
		soundTimer = m.V[x]
	case opcode&0xF0FF == 0xF01E:
		x := (opcode&0x0F00) >> 8
		m.I += uint16(m.V[x])
	case opcode&0xF0FF == 0xF029:
		x := (opcode&0x0F00) >> 8
		digit := m.V[x] & 0x0F 
		m.I = 0x50 + uint16(digit)*5
	case opcode&0xF0FF == 0xF033:
		x := (opcode&0x0F00) >> 8
		value := m.V[x]
		m.memory[m.I] = value / 100
		m.memory[m.I+1] = (value / 10) % 10
		m.memory[m.I+2] = value % 10
	case opcode&0xF0FF == 0xF055:
		x := (opcode & 0x0F00) >> 8
		for i := uint16(0); i <= x; i++ {
			m.memory[m.I+i] = m.V[i]
		}
	case opcode&0xF0FF == 0xF065:
		x := (opcode & 0x0F00) >> 8
		for i := uint16(0); i <= x; i++ {
			m.V[i] = m.memory[m.I+i]
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


func isKeyPressed(chip8Key byte) bool {
    for keyboardKey, mappedKey := range KeyMap_KeyBoard {
        if mappedKey == chip8Key {
            switch keyboardKey {
            case "1": return ebiten.IsKeyPressed(ebiten.Key1)
            case "2": return ebiten.IsKeyPressed(ebiten.Key2)
            case "3": return ebiten.IsKeyPressed(ebiten.Key3)
            case "4": return ebiten.IsKeyPressed(ebiten.Key4)
            case "Q": return ebiten.IsKeyPressed(ebiten.KeyQ)
            case "W": return ebiten.IsKeyPressed(ebiten.KeyW)
            case "E": return ebiten.IsKeyPressed(ebiten.KeyE)
            case "R": return ebiten.IsKeyPressed(ebiten.KeyR)
            case "A": return ebiten.IsKeyPressed(ebiten.KeyA)
            case "S": return ebiten.IsKeyPressed(ebiten.KeyS)
            case "D": return ebiten.IsKeyPressed(ebiten.KeyD)
            case "F": return ebiten.IsKeyPressed(ebiten.KeyF)
            case "Z": return ebiten.IsKeyPressed(ebiten.KeyZ)
            case "X": return ebiten.IsKeyPressed(ebiten.KeyX)
            case "C": return ebiten.IsKeyPressed(ebiten.KeyC)
            case "V": return ebiten.IsKeyPressed(ebiten.KeyV)
            }
        }
    }
    return false
}


func getKeyPressed() byte {
    keyboardKeys := []struct{
        key ebiten.Key
        chip8Key byte
    }{
        {ebiten.Key1, 0x1}, {ebiten.Key2, 0x2}, {ebiten.Key3, 0x3}, {ebiten.Key4, 0xC},
        {ebiten.KeyQ, 0x4}, {ebiten.KeyW, 0x5}, {ebiten.KeyE, 0x6}, {ebiten.KeyR, 0xD},
        {ebiten.KeyA, 0x7}, {ebiten.KeyS, 0x8}, {ebiten.KeyD, 0x9}, {ebiten.KeyF, 0xE},
        {ebiten.KeyZ, 0xA}, {ebiten.KeyX, 0x0}, {ebiten.KeyC, 0xB}, {ebiten.KeyV, 0xF},
    }
    
    for _, keyMap := range keyboardKeys {
        if inpututil.IsKeyJustPressed(keyMap.key) {
            return keyMap.chip8Key
        }
    }
    return 0xFF 
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

func loadTestProgram(m *Memory){
	// Instructions - Display centered, larger IBM logo
	testInstructions := []byte{
		0x00, 0xE0,       // CLS
		0x60, 0x10,       // LD V0, 16 (x position - more centered)
		0x61, 0x0C,       // LD V1, 12 (y position - more centered)
		0xA2, 0x2A,       // LD I, 0x22A (point to "I" sprite)
		0xD0, 0x18,       // DRW V0, V1, 8 (draw 8-byte tall sprite)
		0x70, 0x0A,       // ADD V0, 10 (move x by 10 pixels for spacing)
		0xA2, 0x32,       // LD I, 0x232 (point to "B" sprite)
		0xD0, 0x18,       // DRW V0, V1, 8 (draw 8-byte tall sprite)
		0x70, 0x0A,       // ADD V0, 10 (move x by 10 pixels for spacing)
		0xA2, 0x3A,       // LD I, 0x23A (point to "M" sprite)
		0xD0, 0x18,       // DRW V0, V1, 8 (draw 8-byte tall sprite)
		0x12, 0x1C,       // JP 0x21C (infinite loop to halt)
	}

	// Load program into memory starting at 0x200
	for i, instruction := range testInstructions {
		m.memory[0x200+i] = instruction
	}

	// Larger sprite data for "IBM" - 8 bytes each for better visibility
	ibmSprites := []byte{
		// "I" sprite (8 bytes tall)
		0xFF, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18, 0xFF,
		
		// "B" sprite (8 bytes tall)
		0xFE, 0x33, 0x33, 0xFE, 0x33, 0x33, 0x33, 0xFE,
		
		// "M" sprite (8 bytes tall)
		0x83, 0xC7, 0xEF, 0xDB, 0x83, 0x83, 0x83, 0x83,
	}

	// Load sprites into memory starting at 0x22A
	for i, spriteByte := range ibmSprites {
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
				rect := image.Rect(x*15, y*15, (x+1)*15, (y+1)*15)
				ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), 
					float64(rect.Dx()), float64(rect.Dy()), color.RGBA{255, 255, 255, 255})
			}
		}
	}
	
	
	ebitenutil.DebugPrint(screen, fmt.Sprintf("PC: 0x%03X I: 0x%03X V0: 0x%02X V1: 0x%02X", 
		g.memory.PC, g.memory.I, g.memory.V[0], g.memory.V[1]))
	

	controls := "Controls: 1 2 3 4 | Q W E R | A S D F | Z X C V"
	ebitenutil.DebugPrintAt(screen, controls, 10, 500)
	
	
	keyMapping := "CHIP-8: 1 2 3 C | 4 5 6 D | 7 8 9 E | A 0 B F"
	ebitenutil.DebugPrintAt(screen, keyMapping, 10, 520)
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int){
	return 960, 540
}


