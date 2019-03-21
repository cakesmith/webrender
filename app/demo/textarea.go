package demo

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app/component"
	"github.com/cakesmith/webrender/output"
	"image"
)

var (
	log = logrus.New()
)

func TextArea(bounds image.Rectangle) *component.TextArea {
	taBorder := component.Border{
		output.ColorTerminalGreen,
		1,
	}

	ta := component.NewTextArea(output.ColorBackground, output.ColorTerminalGreen, taBorder, 8, 11, bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)

	ta.Init = func() {

		ta.Buffer = []int{}

		//***** TEST PATTERN *****

		xch := ta.Width() / ta.CharWidth
		ych := ta.Height() / ta.CharHeight

		str := "ALL WORK AND NO PLAY MAKES JACK A DULL BOY. "

		log.Println(str)

		for i := 0; i+len(str) < xch*ych; i = i + len(str) {
			runes := []rune(str)

			for _, r := range runes {
				ta.Buffer = append(ta.Buffer, int(r))
			}
		}

		end := xch*ych - len(ta.Buffer)

		runes := []rune(str[0:end])
		for _, r := range runes {
			ta.Buffer = append(ta.Buffer, int(r))
		}

		//************************

		ta.Draw()

	}

	ta.Draw = func() {

		ta.DrawRectangle(ta.Component.Rectangle, ta.Border.Color)
		ta.DrawRectangle(ta.Component.Rectangle.Inset(ta.Border.Thickness), ta.Component.Color)

		ta.CursorX, ta.CursorY = 0, 0

		for _, chr := range ta.Buffer {
			ta.PrintChar(chr)
		}
	}

	ta.OnKeypress = func(key int) {
		log.WithField("key", key)
		ta.Buffer = append(ta.Buffer, key)
		ta.PrintChar(key)
	}

	return ta

}
