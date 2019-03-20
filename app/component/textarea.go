package component

import (
	"image"
	"image/color"
)

type TextArea struct {
	*Component
	Border
	TextColor             color.Color
	cursorX, cursorY      int
	CharWidth, CharHeight int
	CharMap               *Mapping
}

func NewTextArea(backgroundColor color.Color, textColor color.Color, border Border, charWidth, charHeight, x, y, w, h int) *TextArea {

	buf := []int{}

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

	t.Init = func() {

		chw := t.width() / t.CharWidth
		chh := t.height() / t.CharHeight

		str := "ALL WORK AND NO PLAY MAKES JACK A DULL BOY. "

		log.Println(str)

		for i := 0; i+len(str) < chw*chh; i = i + len(str) {
			t.PrintString(str)
		}

	}

	t.Draw = func() {

		t.DrawRectangle(t.Component.Rectangle, t.Border.Color)
		t.DrawRectangle(t.Component.Rectangle.Inset(t.Border.Thickness), t.Component.Color)

		t.cursorX, t.cursorY = 0, 0

		for _, chr := range buf {
			t.PrintChar(chr)
		}
	}

	t.OnKeypress = func(key int) {
		log.WithField("key", key)
		buf = append(buf, key)
		t.PrintChar(key)
	}

	return t
}

//Bit returns true if the jth bit of x is 1, and false otherwise.
func Bit(x, j uint) bool {
	return !(x&(1<<j) == 0)
}

func (t *TextArea) print(ch int) {

	startX, startY := t.cursorX*t.CharWidth+t.Border.Thickness, t.cursorY*t.CharHeight+t.Border.Thickness

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

func (t *TextArea) height() int {
	return t.Component.Max.Y - t.Component.Min.Y - t.Border.Thickness
}

func (t *TextArea) width() int {
	return t.Component.Max.X - t.Component.Min.X - t.Border.Thickness
}

//Advances the cursor to the beginning of the next line.
func (t *TextArea) println() {
	t.cursorX = 0
	t.cursorY++
}

func (t *TextArea) PrintChar(ch int) {

	if !(t.cursorX < t.width()/t.CharWidth) {
		if t.cursorY < t.height()/t.CharHeight {
			t.println()
		} else {
			t.cursorX, t.cursorY = 0, 0
		}
	}

	t.print(ch)
	t.cursorX++
}

func (t *TextArea) PrintString(str string) {
	runes := []rune(str)

	for _, r := range runes {
		t.PrintChar(int(r))
	}
}
