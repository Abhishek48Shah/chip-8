package internal

import (
	"image/color"
	"math/rand"
	"os"
	"time"
)

type Emulator struct {
	memory [4096]byte
	opcode uint16
	V      [16]uint8
	I      uint16
	Stack  [16]uint16
	SP     uint8
	PC     uint16
}

var chip8Font = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
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
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func NewEmulator(path string) *Emulator {
	emulator := Emulator{}
	data, err := os.ReadFile(path)
	if err != nil {
		panic("error reading file: " + err.Error())
	}
	for i, value := range data {
		emulator.memory[512+i] = value
	}
	for i, value := range chip8Font {
		emulator.memory[i] = value
	}
	emulator.PC = 512
	return &emulator
}
func (e *Emulator) Run(display *Display) {
	for {
		e.opcode = uint16(e.memory[e.PC])<<8 | uint16(e.memory[e.PC+1])
		e.PC += 2
		switch e.opcode & 0xF000 {
		case 0x0000:
			switch e.opcode & 0x00FF {
			case 0x00E0:
				for x := 0; x < 64; x++ {
					for y := 0; y < 32; y++ {
						display.image.Set(x, y, color.Black)
					}
				}
			case 0x00EE:
				e.SP--
				e.PC = e.Stack[e.SP]

			}
		case 0x1000:
			e.PC = e.opcode & 0x0FFF
		case 0x2000:
			NNN := e.opcode & 0x0FFF
			e.Stack[e.SP] = e.PC
			e.SP++
			e.PC = NNN
		case 0x3000:
			NN := e.opcode & 0x00FF
			VX := uint16(e.V[(e.opcode&0x0F00)>>8])
			if VX == NN {
				e.PC += 2
			}
		case 0x4000:
			VX := uint16(e.V[(e.opcode&0x0F00)>>8])
			NN := e.opcode & 0x00FF
			if VX != NN {
				e.PC += 2
			}
		case 0x5000:
			VX := e.V[(e.opcode&0x0F00)>>8]
			VY := e.V[(e.opcode&0x00F0)>>4]
			if VX == VY {
				e.PC += 2
			}
		case 0x6000:
			e.V[(e.opcode&0x0F00)>>8] = byte(e.opcode & 0x00FF)
		case 0x7000:
			e.V[(e.opcode&0x0F00)>>8] += byte(e.opcode & 0x00FF)
		case 0x8000:
			switch e.opcode & 0x000F {
			case 0x0:
				e.V[(e.opcode&0x0F00)>>8] = e.V[(e.opcode&0x00F0)>>4]
			case 0x1:
				e.V[(e.opcode&0x0F00)>>8] |= e.V[(e.opcode&0x00F0)>>4]
			case 0x2:
				e.V[(e.opcode&0x0F00)>>8] &= e.V[(e.opcode&0x00F0)>>4]
			case 0x3:
				e.V[(e.opcode&0x0F00)>>8] ^= e.V[(e.opcode&0x00F0)>>4]
			case 0x4:
				sum := uint16(e.V[(e.opcode&0x0F00)>>8] + e.V[(e.opcode&0x00F0)>>4])
				if sum > 0xFF {
					e.V[0xF] = 1
				} else {
					e.V[0xF] = 0
				}
				e.V[(e.opcode&0x0F00)>>8] = byte(sum)
			case 0x5:
				if e.V[(e.opcode&0x0F00)>>8] >= e.V[(e.opcode&0x00F0)>>4] {
					e.V[0xF] = 1
				} else {
					e.V[0xF] = 0
				}
				e.V[(e.opcode&0x0F00)>>8] -= e.V[(e.opcode&0x00F0)>>4]
			case 0x6:
				X := (e.opcode & 0x0F00) >> 8
				e.V[0xF] = e.V[X] & 0x01
				e.V[X] >>= 1
			case 0x7:
				if e.V[(e.opcode&0x00F0)>>4] >= e.V[(e.opcode&0x0F00)>>8] {
					e.V[0xF] = 1
				} else {
					e.V[0xF] = 0
				}
				e.V[(e.opcode&0x0F00)>>8] = e.V[(e.opcode&0x00F0)>>4] - e.V[(e.opcode&0x0F00)>>8]
			case 0xE:
				if e.V[(e.opcode&0x0F00)>>8]>>7 == 1 {
					e.V[0xF] = 1
				} else {
					e.V[0xF] = 0
				}
				e.V[(e.opcode&0x0F00)>>8] <<= 1
			}
		case 0x9000:
			if e.V[(e.opcode&0x0F00)>>8] != e.V[(e.opcode&0x00F0)>>4] {
				e.PC += 2
			}
		case 0xA000:
			e.I = e.opcode & 0x0FFF
		case 0xB000:
			NNN := e.opcode & 0x0FFF
			e.PC = uint16(e.V[0]) + NNN
		case 0xC000:
			randomNum := uint8(rand.Intn(256))
			NN := uint8(e.opcode & 0x00FF)
			e.V[(e.opcode&0x0F00)>>8] = randomNum & NN
		case 0xD000:
			N := e.opcode & 0x000F
			X := (e.opcode & 0x0F00) >> 8
			Y := (e.opcode & 0x00F0) >> 4
			e.V[0xF] = 0
			for y := 0; y < int(N); y++ {
				sprite := e.memory[int(e.I)+y]
				for x := 0; x < 8; x++ {
					xCoord := (int(e.V[X]) + x) % 64
					yCoord := (int(e.V[Y]) + y) % 32
					if sprite&(0x80>>x) != 0 {
						if display.IsSet(xCoord, yCoord) {
							e.V[0xF] = 1
						}
						display.SetPixel(xCoord, yCoord)
					}
				}
			}
		}
		time.Sleep(1200 * time.Microsecond)
	}
}
