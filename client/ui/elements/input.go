package elements

import (
	"strconv"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

type InputHandler func(string)

type Input struct {
	Length  int
	Cursor  int
	X       int
	Y       int
	Offset  int
	Active  bool
	Handler InputHandler
	message string
}

func NewInput(length int, handler InputHandler) *Input {
	return &Input{length, 0, 0, 0, 0, false, handler, ""}
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
	} else if key == gc.KEY_ENTER || key == gc.KEY_RETURN {
		if input.Handler != nil {
			input.Handler(input.message)
		}
		input.message = ""
		input.Cursor = 0
		input.Offset = 0
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

func (input *Input) Draw(window *gc.Window) {
	printedMessage := input.message[input.Offset:]
	if len(printedMessage) < input.Length {
		printedMessage += strings.Repeat(" ", input.Length-len(printedMessage))
	}
	realCursor := input.Cursor - input.Offset
	printedMessage = printedMessage[:input.Length]
	window.MovePrint(input.Y, input.X, printedMessage)
	if input.Active {
		window.AttrOn(gc.A_REVERSE)
	}
	window.MoveAddChar(input.Y, input.X+realCursor, gc.Char(printedMessage[realCursor]))
	if input.Active {
		window.AttrOff(gc.A_REVERSE)
	}
}

func (input *Input) Move(x, y int) {
	input.X = x
	input.Y = y
}

func (input *Input) Resize(length int) {
	input.Length = length
}
