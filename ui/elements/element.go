package elements

import gc "github.com/mintalk/goncurses"

type Element interface {
	Update(key gc.Key)
	Draw(window *gc.Window)
}
