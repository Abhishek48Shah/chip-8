package main

import (
	"github.com/Abhishek48Shah/chip-8/internal"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	display := internal.NewDisplay()
	emulator := internal.NewEmulator("/home/abhishek/buffer/chip8-roms/programs/Chip8 Picture.ch8")
	go emulator.Run(display)
	err := ebiten.RunGame(display)
	if err != nil {
		panic(err)
	}
}
