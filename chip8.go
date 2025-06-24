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
//should probably put in a struct
var Font_data = []byte{0xF0, 0x90, 0x90, 0x90, 0xF0,
	0x20, 0x60, 0x20, 0x20, 0x70,
	0xF0, 0x10, 0xF0, 0x80, 0xF0,
	0xF0, 0x10, 0xF0, 0x10, 0xF0,
	0x90, 0x90, 0xF0, 0x10, 0x10,
	0xF0, 0x80, 0xF0, 0x10, 0xF0,
	0xF0, 0x80, 0xF0, 0x90, 0xF0,
	0xF0, 0x10, 0x20, 0x40, 0x40,
	0xF0, 0x90, 0xF0, 0x90, 0xF0,
	0xF0, 0x90, 0xF0, 0x10, 0xF0,
	0xF0, 0x90, 0xF0, 0x90, 0x90,
	0xE0, 0x90, 0xE0, 0x90, 0xE0,
	0xF0, 0x80, 0x80, 0x80, 0xF0,
	0xE0, 0x90, 0x90, 0x90, 0xE0,
	0xF0, 0x80, 0xF0, 0x80, 0xF0,
	0xF0, 0x80, 0xF0, 0x80, 0x80}

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
	
	case opcode&0xF0FF == 0xE09E:
		x := (opcode&0x0F00) >> 8
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
		x := (opcode&0x0F00) >> 8
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

var KeyMap_KeyBoard = map[string]byte{
	"1": 0x1, "2": 0x2, "3": 0x3, "4": 0xC,
	"Q": 0x4, "W": 0x5, "E": 0x6, "R": 0xD,
	"A": 0x7, "S": 0x8, "D": 0x9, "F": 0xE,
	"Z": 0xA, "X": 0x0, "C": 0xB, "V": 0xF,
}

var ReverseKeyMap = map[byte]string{
	0x1: "1", 0x2: "2", 0x3: "3", 0xC: "4",
	0x4: "Q", 0x5: "W", 0x6: "E", 0xD: "R",
	0x7: "A", 0x8: "S", 0x9: "D", 0xE: "F",
	0xA: "Z", 0x0: "X", 0xB: "C", 0xF: "V",
}

func isKeyPressed(chip8Key byte) bool {
    keyboardKey, exists := ReverseKeyMap[chip8Key]
    if !exists {
        return false
    }
    
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
var targetInstructionsPerSecond = 500 // 500-700 range from manual

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
	testInstructions := []byte{
		0x00, 0xE0,
		0x60, 0x10,
		0x61, 0x0C,
		0xA2, 0x2A,
		0xD0, 0x18,
		0x70, 0x0A,
		0xA2, 0x32,
		0xD0, 0x18,
		0x70, 0x0A,
		0xA2, 0x3A,
		0xD0, 0x18,
		0x12, 0x1C,
	}

	for i, instruction := range testInstructions {
		m.memory[0x200+i] = instruction
	}

	ibmSprites := []byte{
		0xFF, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18, 0xFF,
		0xFE, 0x33, 0x33, 0xFE, 0x33, 0x33, 0x33, 0xFE,
		0x83, 0xC7, 0xEF, 0xDB, 0x83, 0x83, 0x83, 0x83,
	}

	for i, spriteByte := range ibmSprites {
		m.memory[0x22A+i] = spriteByte
	}
}

func loadROM(m *Memory, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %v", err)
	}
	
	for i, byte := range data {
		if 0x200+i >= 4096 {
			break 
		}
		m.memory[0x200+i] = byte
	}
	
	return nil
}

func(g *Game) Draw(screen *ebiten.Image){
	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			if display_grid[x][y] {
				rect := image.Rect(x*8, y*8, (x+1)*8, (y+1)*8)
				ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), 
					float64(rect.Dx()), float64(rect.Dy()), color.RGBA{255, 255, 255, 255})
			}
		}
	}
	
	// Debug info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("PC: 0x%03X I: 0x%03X", g.memory.PC, g.memory.I))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("V0: 0x%02X V1: 0x%02X V2: 0x%02X V3: 0x%02X", 
		g.memory.V[0], g.memory.V[1], g.memory.V[2], g.memory.V[3]), 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DT: %d ST: %d", delayTimer, soundTimer), 10, 50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Target IPS: %d | Current: %d", targetInstructionsPerSecond, g.instructionsPerFrame*60), 10, 70)
	
	// Simple key mapping reference
	ebitenutil.DebugPrintAt(screen, "Keys: 1234 QWER ASDF ZXCV -> 123C 456D 789E A0BF", 10, 290)
	ebitenutil.DebugPrintAt(screen, "Speed: +/- to adjust emulation speed", 10, 310)
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int){
	return 960, 540
}

type Game struct {
	memory *Memory
	stack  *Stack
	lastUpdate time.Time
	instructionsPerFrame int
}

func(g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		targetInstructionsPerSecond += 100
		if targetInstructionsPerSecond > 2000 {
			targetInstructionsPerSecond = 2000
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		targetInstructionsPerSecond -= 100
		if targetInstructionsPerSecond < 100 {
			targetInstructionsPerSecond = 100
		}
	}
	now := time.Now()
	deltaTime := now.Sub(g.lastUpdate)
	g.lastUpdate = now
	g.instructionsPerFrame = int(float64(targetInstructionsPerSecond) * deltaTime.Seconds())
		if g.instructionsPerFrame < 1 {
		g.instructionsPerFrame = 1
	}

	for i := 0; i < g.instructionsPerFrame; i++ {
		opcode := fetch(g.memory)
		execute(opcode, g.memory, g.stack)
	}
	
	return nil
}

