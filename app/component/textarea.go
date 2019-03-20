package component

import (
	"image/color"
)

type TextArea struct {
	Component
	Border
	color.Color
	Background            color.Color
	cursorX, cursorY      int
	CharWidth, CharHeight int
	CharMap               *Mapping
}


//Bit returns true if the jth bit of x is 1, and false otherwise.
func Bit(x, j int) bool {
	return !(x&(1<<uint(j)) == 0)
}

//func (t *TextArea) print(ch int) {
//
//	startX, startY := t.cursorX*t.CharWidth, t.cursorY*t.CharHeight
//
//	stopX, stopY := startX+t.CharWidth, startY+t.CharHeight
//
//	printMe := t.CharMap.Get(ch)
//
//	for y := int(0); startY+y < stopY; y++ {
//		for x := int(0); startX+x < stopX; x++ {
//
//			color := t.Color
//
//			pixel := Bit(printMe[y], x)
//
//			if !pixel {
//				color = t.Background
//			}
//
//			t.DrawPixel(startX+x, startY+y, color)
//
//		}
//	}
//
//}
//
////Advances the cursor to the beginning of the next line.
//func (t *TextArea) println() {
//	if t.cursorY < t.Height/t.CharHeight {
//		t.cursorX = 0
//		t.cursorY++
//	}
//}
//
//func (t *TextArea) PrintChar(ch int) {
//
//	if t.cursorX < t.Width/t.CharWidth {
//
//		t.print(ch)
//		t.cursorX++
//
//	} else {
//		if t.cursorY < t.Height/t.CharHeight {
//			t.cursorX = 0
//			t.cursorY = 0
//			t.print(ch)
//			t.cursorX = 1
//		} else {
//			t.println()
//		}
//	}
//}
