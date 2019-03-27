package font

import (
	"bytes"
	"fmt"
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

	if len(buf.Bytes()) == 0 {
		fmt.Println("expected draw function to be called")
	}

	main.OnClick(0, opts.X + opts.CellWidth * 2 -1, opts.Y + opts.CellHeight * 2 -1)

	expected := "0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0"
	actual := grid.String()

	if expected != actual {
		t.Errorf("expected\n%v\nreceived\n%v", expected, actual)
	}



}
