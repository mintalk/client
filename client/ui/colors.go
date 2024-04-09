package ui

import gc "github.com/rthornton128/goncurses"

func InitColors() {
	gc.StartColor()
	gc.InitPair(1, gc.C_BLUE, gc.C_BLACK)  // Active panel border
	gc.InitPair(2, gc.C_WHITE, gc.C_BLACK) // Inactive panel border
}
