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

	b.Component.Draw = func() {
		b.Draw()
	}

	b.Component.Init = func(){}

	return &b
}

func (b *Button) Draw() {
	b.DrawRectangle(b.Bounds(), b.Border)
	b.DrawRectangle(b.Bounds().Inset(b.Border.Thickness), b.Component.Color)
}
