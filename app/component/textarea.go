package component

import (
	"image"
	"image/color"
	"math"
)

type TextArea struct {
	*Component
	Border
	TextColor             color.Color
	CursorX, CursorY      int
	CharWidth, CharHeight int
	CharMap               *Mapping
	Buffer                []int
}

func NewTextArea(backgroundColor color.Color, textColor color.Color, border Border, charWidth, charHeight, x, y, w, h int) *TextArea {

	t := &TextArea{
		Component: &Component{
			Rectangle: image.Rect(x, y, x+w, y+h),
			Color:     backgroundColor,
			Init:      nil,
			OnClick:   nil,
		},
		Border:     border,
		TextColor:  textColor,
		CharWidth:  charWidth,
		CharHeight: charHeight,
		CharMap:    NewMapping(),
	}

	return t
}

//Bit returns true if the jth bit of x is 1, and false otherwise.
func Bit(x, j uint) bool {
	return !(x&(1<<j) == 0)
}

func (t *TextArea) _print(ch int) {

	//TODO determine where this padding should live
	// because it's part of the calculations here, it
	// can change the value of the true height/width of
	// this component.
	//

	padX := int(math.Mod(float64(t.Width()), float64(t.CharWidth)) / 2)
	padY := int(math.Mod(float64(t.Height()), float64(t.CharHeight)) / 2)

	startX, startY := t.CursorX*t.CharWidth+t.Border.Thickness+padX, t.CursorY*t.CharHeight+t.Border.Thickness+padY

	stopX, stopY := startX+t.CharWidth, startY+t.CharHeight

	printMe := t.CharMap.Get(ch)

	for y := int(0); startY+y < stopY; y++ {
		for x := int(0); startX+x < stopX; x++ {

			c := t.TextColor

			pixel := Bit(uint(printMe[y]), uint(x))

			if !pixel {
				c = t.Component.Color
			}

			t.DrawPixel(startX+x, startY+y, c)

		}
	}

}

func (t *TextArea) Height() int {
	return t.Component.Max.Y - t.Component.Min.Y - t.Border.Thickness
}

func (t *TextArea) Width() int {
	return t.Component.Max.X - t.Component.Min.X - t.Border.Thickness
}

//Advances the cursor to the beginning of the next line.
func (t *TextArea) PrintLn() {
	t.CursorX = 0
	t.CursorY++
}

func (t *TextArea) PrintChar(ch int) {

	if !(t.CursorX < t.Width()/t.CharWidth) {
		if t.CursorY < t.Height()/t.CharHeight {
			t.PrintLn()
		} else {
			t.CursorX, t.CursorY = 0, 0
		}
	}

	t._print(ch)
	t.CursorX++
}

func (t *TextArea) PrintString(str string) {
	runes := []rune(str)

	for _, r := range runes {
		t.PrintChar(int(r))
	}
}
