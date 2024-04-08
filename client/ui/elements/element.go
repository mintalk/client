package elements

import gc "github.com/rthornton128/goncurses"

type Element interface {
	Update(key gc.Key)
	Draw(window *gc.Window, x, y int)
}
