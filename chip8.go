package main

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
