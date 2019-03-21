package font

import (
	"bytes"
	"github.com/cakesmith/webrender/app/component"
	"github.com/cakesmith/webrender/output"
	"image"
	"testing"
)

var width = 640
var height = 480

func createContainer(w, h int) (*component.Container, *bytes.Buffer) {

	buf := new(bytes.Buffer)
	main := component.Container{
		Rectangle: image.Rect(0, 0, w, h),
		Terminal:  output.Terminal{Writer: buf},
	}
	return &main, buf
}

func TestCell(t *testing.T) {

	main, buf := createContainer(width, height)

	background := output.ColorRed

	cell := NewCell(image.Rect(0, 0, 100, 100), background, output.ColorTerminalGreen)

	main.Add(cell.Component)

	expected := "r 0 0 100 100 200:0:0:0"

	actual := string(buf.Bytes())

	if expected != actual {
		t.Errorf("expected %v received %v", expected, actual)
	}

	if cell.active {
		t.Error("expected cell.active to be false")
	}

	buf.Reset()

	cell.OnClick(0, 1, 1)

	expected = "r 0 0 100 100 51:255:51:0"

	actual = string(buf.Bytes())

	if expected != actual {
		t.Errorf("expected %v received %v", expected, actual)
	}

	if !cell.active {
		t.Error("expected cell.active to be true")
	}


}

func TestGrid(t *testing.T) {

	main, buf := createContainer(width, height)

	opts := GridOptions{
		CellWidth:       16,
		CellHeight:      22,
		XCells:          8,
		YCells:          11,
		X:               120,
		Y:               240,
		Thickness:       1,
		BackgroundColor: output.ColorBackground,
		LineColor:       output.ColorBlack,
		ActiveColor:     output.ColorTerminalGreen,
	}

	grid := NewGrid(opts)

	main.Add(grid.Component)

	expected := "r 0 0 16 22 40:40:40:0r 16 0 16 22 40:40:40:0r 32 0 16 22 40:40:40:0r 48 0 16 22 40:40:40:0r 64 0 16 22 40:40:40:0r 80 0 16 22 40:40:40:0r 96 0 16 22 40:40:40:0r 112 0 16 22 40:40:40:0r 0 22 16 22 40:40:40:0r 16 22 16 22 40:40:40:0r 32 22 16 22 40:40:40:0r 48 22 16 22 40:40:40:0r 64 22 16 22 40:40:40:0r 80 22 16 22 40:40:40:0r 96 22 16 22 40:40:40:0r 112 22 16 22 40:40:40:0r 0 44 16 22 40:40:40:0r 16 44 16 22 40:40:40:0r 32 44 16 22 40:40:40:0r 48 44 16 22 40:40:40:0r 64 44 16 22 40:40:40:0r 80 44 16 22 40:40:40:0r 96 44 16 22 40:40:40:0r 112 44 16 22 40:40:40:0r 0 66 16 22 40:40:40:0r 16 66 16 22 40:40:40:0r 32 66 16 22 40:40:40:0r 48 66 16 22 40:40:40:0r 64 66 16 22 40:40:40:0r 80 66 16 22 40:40:40:0r 96 66 16 22 40:40:40:0r 112 66 16 22 40:40:40:0r 0 88 16 22 40:40:40:0r 16 88 16 22 40:40:40:0r 32 88 16 22 40:40:40:0r 48 88 16 22 40:40:40:0r 64 88 16 22 40:40:40:0r 80 88 16 22 40:40:40:0r 96 88 16 22 40:40:40:0r 112 88 16 22 40:40:40:0r 0 110 16 22 40:40:40:0r 16 110 16 22 40:40:40:0r 32 110 16 22 40:40:40:0r 48 110 16 22 40:40:40:0r 64 110 16 22 40:40:40:0r 80 110 16 22 40:40:40:0r 96 110 16 22 40:40:40:0r 112 110 16 22 40:40:40:0r 0 132 16 22 40:40:40:0r 16 132 16 22 40:40:40:0r 32 132 16 22 40:40:40:0r 48 132 16 22 40:40:40:0r 64 132 16 22 40:40:40:0r 80 132 16 22 40:40:40:0r 96 132 16 22 40:40:40:0r 112 132 16 22 40:40:40:0r 0 154 16 22 40:40:40:0r 16 154 16 22 40:40:40:0r 32 154 16 22 40:40:40:0r 48 154 16 22 40:40:40:0r 64 154 16 22 40:40:40:0r 80 154 16 22 40:40:40:0r 96 154 16 22 40:40:40:0r 112 154 16 22 40:40:40:0r 0 176 16 22 40:40:40:0r 16 176 16 22 40:40:40:0r 32 176 16 22 40:40:40:0r 48 176 16 22 40:40:40:0r 64 176 16 22 40:40:40:0r 80 176 16 22 40:40:40:0r 96 176 16 22 40:40:40:0r 112 176 16 22 40:40:40:0r 0 198 16 22 40:40:40:0r 16 198 16 22 40:40:40:0r 32 198 16 22 40:40:40:0r 48 198 16 22 40:40:40:0r 64 198 16 22 40:40:40:0r 80 198 16 22 40:40:40:0r 96 198 16 22 40:40:40:0r 112 198 16 22 40:40:40:0r 0 220 16 22 40:40:40:0r 16 220 16 22 40:40:40:0r 32 220 16 22 40:40:40:0r 48 220 16 22 40:40:40:0r 64 220 16 22 40:40:40:0r 80 220 16 22 40:40:40:0r 96 220 16 22 40:40:40:0r 112 220 16 22 40:40:40:0"
	actual := string(buf.Bytes())

	if expected != actual {
		t.Errorf("expected\n%v\nreceived\n%v", expected, actual)
	}

}
