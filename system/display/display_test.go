package display

import (
	"bytes"
	"testing"
)

func TestDisplayWriter_DrawPixel(t *testing.T) {

	buf := new(bytes.Buffer)

	display := Terminal{buf}

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
