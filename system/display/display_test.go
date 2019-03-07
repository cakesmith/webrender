package display

import (
	"bytes"
	"testing"
)

func Test_MakePacket(t *testing.T) {

}

func Test_Send(t *testing.T) {

}

func Test_DrawPixel(t *testing.T) {

	buf := new(bytes.Buffer)

	display := Terminal{
		Writer: buf,
	}

	err := display.DrawPixel(100, 200, ColorBlack)
	if err != nil {
		t.Error(err)
	}

	expected := "r 100 200 1 1 " + ColorBlack.String()

	if expected != buf.String() {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}

	buf.Reset()

	err = display.DrawPixel(100, 200, ColorWhite)
	if err != nil {
		t.Error(err)
	}

	expected = "r 100 200 1 1 " + ColorWhite.String()

	if expected != buf.String() {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}

}

func Test_DrawVert(t *testing.T) {

}

func Test_DrawHoriz(t *testing.T) {

}

func Test_DrawLine(t *testing.T) {

}

func TestAbs(t *testing.T) {

}

func TestSqrt(t *testing.T) {

}

func TestDrawCircle(t *testing.T) {

}
