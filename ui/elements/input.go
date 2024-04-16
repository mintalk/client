package elements

import (
	"strings"
	"unicode/utf8"

	gc "github.com/mintalk/goncurses"
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
	message []byte
}

func NewInput(length int, handler InputHandler) *Input {
	return &Input{length, 0, 0, 0, 0, false, handler, []byte{}}
}

func (input *Input) Update(key gc.Key) {
	if !input.Active {
		return
	}
	if key == 0 {
		return
	}
	if key == gc.KEY_BACKSPACE {
		input.Cursor = input.StartCursor()
		if input.Cursor > 0 {
			begin := input.Cursor
			input.Cursor--
			input.Cursor = input.StartCursor()
			end := input.Cursor
			input.message = append(input.message[:end], input.message[begin:]...)
		}
	} else if key == gc.KEY_LEFT {
		if input.StartCursor() > 0 {
			originalCursor := input.StartCursor()
			for input.StartCursor() == originalCursor {
				input.Cursor--
			}
			input.Cursor = input.StartCursor()
		}
	} else if key == gc.KEY_RIGHT {
		if input.Cursor < len(input.message) {
			originalCursor := input.StartCursor()
			for input.StartCursor() == originalCursor {
				input.Cursor++
			}
		}
	} else if key == gc.KEY_ENTER || key == gc.KEY_RETURN {
		input.Handler(string(input.message))
		input.Cursor = 0
		input.message = []byte{}
	} else if key == gc.KEY_DC {
		input.Cursor = input.StartCursor()
		if input.RuneCursor() < input.RuneLength() {
			end := input.Cursor
			for input.StartCursor() == end {
				input.Cursor++
			}
			begin := input.Cursor
			input.message = append(input.message[:end], input.message[begin:]...)
			input.Cursor = end
		}
	} else if key == gc.KEY_HOME {
	} else if key == gc.KEY_END {
	} else {
		input.Cursor = input.StartCursor()
		runeKey := rune(key)
		if utf8.ValidRune(runeKey) && input.RuneLength() < input.Length {
			begin := input.Cursor
			input.message = append(input.message[:begin], append([]byte{byte(key)}, input.message[begin:]...)...)
			input.Cursor++
		}
	}
}

func (input *Input) RuneCursor() int {
	if len(input.message) == 0 {
		return 0
	}
	if input.Cursor >= len(input.message) {
		return input.RuneLength()
	}
	return utf8.RuneCount(input.message[:input.StartCursor()])
}

func (input *Input) RuneLength() int {
	return utf8.RuneCount(input.message)
}

func (input *Input) StartCursor() int {
	if len(input.message) == 0 {
		return 0
	}
	if input.Cursor >= len(input.message) {
		return len(input.message)
	}
	startCursor := input.Cursor
	for !utf8.RuneStart(input.message[startCursor]) {
		startCursor--
	}
	return startCursor
}

func (input *Input) Draw(window *gc.Window) {
	if !input.Active {
		return
	}
	window.MovePrint(input.Y, input.X, strings.Repeat(" ", input.Length))
	window.MovePrint(input.Y, input.X, string(input.message))
	window.AttrOn(gc.A_REVERSE)
	startCursor := input.StartCursor()
	cursorRune := rune(' ')
	if startCursor < len(input.message) {
		cursorRune, _ = utf8.DecodeRune(input.message[startCursor:])
	}
	window.MovePrint(input.Y, input.X+input.RuneCursor(), string(cursorRune))
	window.AttrOff(gc.A_REVERSE)
	window.Refresh()
}

func (input *Input) Move(x, y int) {
	input.X = x
	input.Y = y
}

func (input *Input) Resize(length int) {
	input.Length = length
}
