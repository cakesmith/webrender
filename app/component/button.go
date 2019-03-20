package component

import (
	"image"
	"image/color"
)

type Button struct {
	*Component
	Border
}

func NewButton(color color.Color, border Border, x, y, width, height int) *Button {
	b := Button{
		Component: &Component{
			Color:     color,
			Rectangle: image.Rect(x, y, x+width, y+height),
		},
		Border: border,
	}
	return &b
}

func (b *Button) Draw() {

	x := b.Bounds().Min.X
	y := b.Bounds().Min.Y
	x1 := b.Bounds().Max.X
	y1 := b.Bounds().Max.Y

	b.DrawRectangle(x, y, x1, y1, b.Border)

	bx := x + b.Thickness
	by := y + b.Thickness
	bx1 := x1 - b.Thickness
	by1 := y1 - b.Thickness

	b.DrawRectangle(bx, by, bx1, by1, b.Component.Color)
}
