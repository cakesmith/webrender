package os

import (
	"bytes"
	"testing"
)

func TestDisplayWriter_Send(t *testing.T) {

	buf := new(bytes.Buffer)

	display := DisplayWriter{
		Writer: buf,
	}

	err := display.DrawPixel(100, 200, Pixel{0})
	if err != nil {
		t.Error(err)
	}

	expected := "d 100 200 0"

	if expected != buf.String() {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}

	buf.Reset()

	err = display.DrawPixel(100, 200, Pixel{0})
	if err != nil {
		t.Error(err)
	}

	expected = ""

	if expected != buf.String() {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}

	err = display.DrawPixel(100, 200, Pixel{1})

	expected = "d 100 200 1"

	if expected != buf.String() {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}

	buf.Reset()
	expected = ""

	err = display.DrawPixel(100, 200, Pixel{1})

	if expected != buf.String() {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}


}
