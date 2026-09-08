package internal

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Display struct {
	image *image.RGBA
	key   []uint16
}

func NewDisplay() *Display {
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("Chip-8")
	return &Display{image: image.NewRGBA(image.Rect(0, 0, 64, 32))}
}
func (d *Display) Update() error {
	return nil
}
func (d *Display) Draw(screen *ebiten.Image) {
	screen.WritePixels(d.image.Pix)
}
func (d *Display) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}
func (d *Display) SetPixel(x, y int) {
	d.image.Set(x, y, color.White)
}
func (d *Display) IsSet(x, y int) bool {
	return d.image.At(x, y) == color.White
}
