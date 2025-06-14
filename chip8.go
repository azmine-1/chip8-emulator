package main

import "os"

type Chip8 struct {
	memory     [4096]byte
	V          [16]byte
	I          uint16
	pc         uint16
	gfx        [64 * 32]byte
	delayTimer byte
	soundTimer byte
	stack      [16]uint16
	sp         uint16
	key        [16]byte
	drawFlag   bool
}

func (chip *Chip8) initialize() {
	chip.pc = 0x200
	chip.I = 0
	chip.sp = 0
	chip.drawFlag = false
}

func (chip *Chip8) loadROM(filename string) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		return err
	}
	copy(chip.memory[0x200:], data)
	return nil
}

func (chip *Chip8) EmulateCycle() {
	opcode := uint16(chip.memory[chip.pc])<<8 | uint16(chip.memory[chip.pc+1])
	chip.executeOpcode(opcode)
	chip.updateTimer()

}
