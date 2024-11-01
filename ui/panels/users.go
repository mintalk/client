package panels

import (
	"mintalk/client/cache"
	"strings"

	"github.com/rivo/tview"
)

type Users struct {
	*tview.TextView
	app *tview.Application
	serverCache  *cache.ServerCache
}

func NewUsers(app *tview.Application, serverCache *cache.ServerCache) *Users {
	users := &Users{
		TextView: tview.NewTextView().SetScrollable(true),
		app: app,
		serverCache: serverCache,
	}
	users.SetBorder(true)
	users.SetTitle("Users")
	users.serverCache.AddListener(users.updateListData)
	return users
}

func (users *Users) updateListData() {
	userList := make([]string, 0)
	for _, user := range users.serverCache.Users {
		userList = append(userList, user)
	}
	text := strings.Join(userList, "\n")
	users.app.QueueUpdateDraw(func () {
		users.SetText(text)
	})
}
