package component_test

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
		Rectangle:  image.Rect(0, 0, width, height),
		Terminal:   output.Terminal{Writer: buf},
		Components: []*component.Component{},
	}
	return &main, buf
}

var redBorder = component.Border{
	Color:     output.ColorRed,
	Thickness: 1,
}

func TestFocus(t *testing.T) {

	// scenarios:
	//  - multiple components, no focus set
	//  - multiple components, focus selected
	//  - multiple components, focus changed
	//  - single component, no focus set
	//  - single component, focus selected

	main, _ := createContainer(640, 480)

	testButton := component.NewButton(output.ColorTerminalGreen, redBorder, 100, 200, 50, 75)

	main.Add(testButton.Component)

}

func TestButton(t *testing.T) {

	main, buf := createContainer(640, 480)

	testButton := component.NewButton(output.ColorTerminalGreen, redBorder, 100, 200, 50, 75)

	main.Add(testButton.Component)

	var (
		initCalled               = false
		keyPressed, tbtn, tx, ty int
	)

	testButton.Component.Init = func() {
		initCalled = true
		testButton.Draw()
	}

	testButton.Component.OnKeypress = func(key int) {
		keyPressed = key
	}

	testButton.Component.OnClick = func(btn, x, y int) {
		tbtn, tx, ty = btn, x, y
	}

	main.Init()

	if !initCalled {
		t.Error("init not called")
	}

	cmdstr := string(buf.Bytes())

	// init should be called

	// draw should be called twice
	//
	//r 100 200 150 275 200:0:0:0
	//r 101 201 149 274 51:255:51:0

	expected := "r 100 200 150 275 200:0:0:0r 101 201 149 274 51:255:51:0"
	if cmdstr != expected {
		t.Errorf("expected\n%v\nreceived\n%v", expected, cmdstr)
	}

	main.OnClick(1, 50, 75)

	// click should not have happened

	if tbtn != 0 || tx != 0 || ty != 0 {
		t.Errorf("click should not fire\nbtn: %v x: %v y: %v", tbtn, tx, ty)
	}

	main.OnClick(1, 120, 250)

	// click should be (1, 20, 50)

	if tbtn != 1 || tx != 20 || ty != 50 {
		t.Errorf("incorrect click coords\nbtn: %v x: %v y: %v", tbtn, tx, ty)
	}

	// button should have focus since it is the only component
	// and receives keypresses

	main.OnKeypress(55) // the number 7

	if keyPressed != 55 {
		t.Errorf("expected 55 received %v", keyPressed)
	}

}
