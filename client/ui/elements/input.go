package elements

import (
	"strconv"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

type Input struct {
	Length     int
	Cursor     int
	MarkCursor int
	Offset     int
	Active     bool
	message    string
	mark       bool
}

func NewInput(length int) *Input {
	return &Input{length, 0, 0, 0, false, "", false}
}

func (input *Input) Update(key gc.Key) {
	if !input.Active {
		return
	}
	if key == gc.KEY_BACKSPACE {
		if input.Cursor > 0 {
			input.moveLeft(true)
			newMessage := input.message[:input.Cursor]
			if input.Cursor < len(input.message) {
				newMessage += input.message[input.Cursor+1:]
			}
			input.message = newMessage
		}
	} else if key == gc.KEY_LEFT {
		input.moveLeft(false)
	} else if key == gc.KEY_RIGHT {
		if input.Cursor < len(input.message) {
			input.moveRight()
		}
	} else {
		if !strconv.IsPrint(rune(key)) {
			return
		}
		newMessage := ""
		if input.Cursor > 0 {
			newMessage = input.message[:input.Cursor]
		}
		newMessage += string(rune(key))
		if input.Cursor < len(input.message) {
			newMessage += input.message[input.Cursor:]
		}
		input.message = newMessage
		input.moveRight()
	}
	input.mark = false
}

func (input *Input) moveLeft(scroll bool) {
	if input.Cursor <= 0 {
		return
	}
	input.Cursor--
	if scroll && input.Offset > 0 {
		input.Offset--
	}
	if input.Cursor-input.Offset < 0 {
		input.Offset = input.Cursor
	}
}

func (input *Input) moveRight() {
	if input.Cursor >= len(input.message) {
		return
	}
	input.Cursor++
	if input.Cursor-input.Offset >= input.Length {
		input.Offset++
	}
}

func (input *Input) Draw(window *gc.Window, x, y int) {
	printedMessage := input.message[input.Offset:]
	if len(printedMessage) < input.Length {
		printedMessage += strings.Repeat(" ", input.Length-len(printedMessage))
	}
	realCursor := input.Cursor - input.Offset
	printedMessage = printedMessage[:input.Length]
	window.MovePrint(y, x, printedMessage)
	window.AttrOn(gc.A_REVERSE)
	window.MoveAddChar(y, x+realCursor, gc.Char(printedMessage[realCursor]))
	window.AttrOff(gc.A_REVERSE)
}
