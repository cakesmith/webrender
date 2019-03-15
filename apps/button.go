package apps

import (
	"github.com/cakesmith/webrender/output/color"
)

type Button struct {
	Component
	Border          color.Color
	BorderThickness uint
	Text            Text
}

func (b Button) Hit(mx, my uint) bool {
	return mx > b.X && my > b.Y && mx < b.X+b.Height && my < b.Y+b.Width
}

func (b Button) Draw() {
	b.DrawRectangle(b.X, b.Y, b.Width, b.Height, b.Border)

	bx := b.X + b.BorderThickness
	by := b.Y + b.BorderThickness
	bw := b.Width - b.BorderThickness - 1
	bh := b.Height - b.BorderThickness - 1

	b.DrawRectangle(bx, by, bw, bh, b.Color)
}
