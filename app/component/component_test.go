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
		Rectangle: image.Rect(0, 0, w, h),
		Terminal:  output.Terminal{Writer: buf},
	}
	return &main, buf
}

var redBorder = component.Border{
	Color:     output.ColorRed,
	Thickness: 1,
}

var blackBorder = component.Border{
	Color:     output.ColorBlack,
	Thickness: 2,
}

// focus scenarios:
//  - single component, default focus
//     component should have focus
//  - single component, focus selected
//     component should have focus
//
//  - multiple components, default focus
//      all components have focus
//  - multiple components, single selected
//      selected component should have focus
//  - multiple components, multiple selected
//      selected components should have focus

func TestSingleComponentDefaultFocus(t *testing.T) {

	main, _ := createContainer(width, height)

	testButton := component.NewButton(
		output.ColorTerminalGreen,
		redBorder,
		100,
		200,
		50,
		75,
	)

	main.Add(testButton.Component)

	var keyPressed int

	testButton.Component.OnKeypress = func(key int) {
		keyPressed = key
	}

	main.OnKeypress(56)

	if keyPressed != 56 {
		t.Errorf("expected 56 received %v", keyPressed)
	}

}

func TestSingleComponentSetFocus(t *testing.T) {

	main, _ := createContainer(width, height)

	testButton := component.NewButton(
		output.ColorTerminalGreen,
		redBorder,
		100,
		200,
		50,
		75,
	)

	main.Add(testButton.Component)

	var keyPressed int

	testButton.Component.OnKeypress = func(key int) {
		keyPressed = key
	}

	main.Focused = []*component.Component{
		testButton.Component,
	}

	main.OnKeypress(56)

	if keyPressed != 56 {
		t.Errorf("expected 56 received %v", keyPressed)
	}

}

func TestMultipleComponentsDefaultFocus(t *testing.T) {

	main, _ := createContainer(width, height)

	testButton1 := component.NewButton(
		output.ColorTerminalGreen,
		redBorder,
		100,
		200,
		50,
		75,
	)

	testButton2 := component.NewButton(
		output.ColorRed,
		blackBorder,
		200,
		300,
		50,
		75,
	)

	main.Add(testButton1.Component)
	main.Add(testButton2.Component)

	var keyPressed1, keyPressed2 int

	testButton1.Component.OnKeypress = func(key int) {
		keyPressed1 = key
	}

	testButton2.Component.OnKeypress = func(key int) {
		keyPressed2 = key
	}

	main.OnKeypress(56)

	if keyPressed1 != 56 {
		t.Errorf("keypress 1 expected 56 received %v", keyPressed1)
	}

	if keyPressed2 != 56 {
		t.Errorf("keypress 2 expected 56 received %v", keyPressed1)
	}

}

func TestMultipleComponentsSingleFocus(t *testing.T) {

	main, _ := createContainer(width, height)

	testButton1 := component.NewButton(
		output.ColorTerminalGreen,
		redBorder,
		100,
		200,
		50,
		75,
	)

	testButton2 := component.NewButton(
		output.ColorRed,
		blackBorder,
		200,
		300,
		50,
		75,
	)

	main.Add(testButton1.Component)
	main.Add(testButton2.Component)

	var keyPressed1, keyPressed2 int

	testButton1.Component.OnKeypress = func(key int) {
		keyPressed1 = key
	}

	testButton2.Component.OnKeypress = func(key int) {
		keyPressed2 = key
	}

	main.Focused = []*component.Component{
		testButton1.Component,
	}

	main.OnKeypress(56)

	if keyPressed1 != 56 {
		t.Errorf("keypress 1 expected 56 received %v", keyPressed1)
	}

	if keyPressed2 != 0 {
		t.Errorf("keypress 2 expected 0 received %v", keyPressed1)
	}

}

func TestMultipleComponentsMultipleFocus(t *testing.T) {

	main, _ := createContainer(width, height)

	testButton1 := component.NewButton(
		output.ColorTerminalGreen,
		redBorder,
		100,
		200,
		50,
		75,
	)

	testButton2 := component.NewButton(
		output.ColorRed,
		blackBorder,
		200,
		300,
		50,
		75,
	)

	testButton3 := component.NewButton(
		output.ColorBackground,
		component.Border{},
		251,
		400,
		50,
		75,
	)

	main.Add(testButton1.Component)
	main.Add(testButton2.Component)
	main.Add(testButton3.Component)

	var keyPressed1, keyPressed2, keyPressed3 int

	testButton1.Component.OnKeypress = func(key int) {
		keyPressed1 = key
	}

	testButton2.Component.OnKeypress = func(key int) {
		keyPressed2 = key
	}

	testButton3.Component.OnKeypress = func(key int) {
		keyPressed3 = key
	}

	main.Focused = []*component.Component{
		testButton1.Component,
		testButton3.Component,
	}

	main.OnKeypress(56)

	if keyPressed1 != 56 {
		t.Errorf("keypress 1 expected 56 received %v", keyPressed1)
	}

	if keyPressed2 != 0 {
		t.Errorf("keypress 2 expected 0 received %v", keyPressed2)
	}

	if keyPressed3 != 56 {
		t.Errorf("keypress 3 expected 56 received %v", keyPressed3)
	}
}

func TestNoKeyHandler(t *testing.T) {
	main, _ := createContainer(width, height)

	testButton := component.NewButton(
		output.ColorTerminalGreen,
		redBorder,
		100,
		200,
		50,
		75,
	)

	main.Add(testButton.Component)

	main.Focused = []*component.Component{
		testButton.Component,
	}

	main.OnKeypress(58)

}

func TestMultipleButtons(t *testing.T) {
	main, _ := createContainer(width, height)

	testButton1 := component.NewButton(output.ColorTerminalGreen, redBorder, 100, 200, 50, 75)
	testButton2 := component.NewButton(output.ColorRed, blackBorder, 200, 300, 50, 75)

	main.Add(testButton1.Component)
	main.Add(testButton2.Component)

	var (
		tbtn1, tx1, ty1,
		tbtn2, tx2, ty2 int
	)

	rst := func() {
		tbtn1, tx1, ty1, tbtn2, tx2, ty2 = 0, 0, 0, 0, 0, 0
	}

	chk := func(a, b, c, d, e, f int) {
		t.Helper()
		if tbtn1 != a || tx1 != b || ty1 != c || tbtn2 != d || tx2 != e || ty2 != f {
			t.Errorf("one: %v %v %v\ntwo: %v %v %v\n", tbtn1, tx1, ty1, tbtn2, tx2, ty2)
		}
	}

	testButton1.Component.OnClick = func(btn, x, y int) {
		tbtn1, tx1, ty1 = btn, x, y
	}

	testButton2.Component.OnClick = func(btn, x, y int) {
		tbtn2, tx2, ty2 = btn, x, y
	}

	main.OnClick(1, 125, 221)
	chk(1, 25, 21, 0, 0, 0)
	rst()

	main.OnClick(1, 225, 333)
	chk(0, 0, 0, 1, 25, 33)
	rst()

	main.OnClick(2, 580, 380)
	chk(0, 0, 0, 0, 0, 0)

}

func TestButton(t *testing.T) {

	main, buf := createContainer(width, height)

	testButton := component.NewButton(output.ColorTerminalGreen, redBorder, 100, 200, 50, 75)

	main.Add(testButton.Component)

	cmdstr := string(buf.Bytes())

	// init should be called

	// draw should be called twice

	expected := "r 100 200 50 75 200:0:0:0r 101 201 48 73 51:255:51:0"
	if cmdstr != expected {
		t.Errorf("expected\n%v\nreceived\n%v", expected, cmdstr)
	}

	var (
		initCalled               = false
		keyPressed, tbtn, tx, ty int
	)

	testButton.Component.Init = func() {
		initCalled = true
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

	main.OnKeypress(55)

	if keyPressed != 55 {
		t.Errorf("expected 55 received %v", keyPressed)
	}

}
