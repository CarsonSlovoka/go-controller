package funcs

import (
	"fmt"
	"github.com/go-vgo/robotgo"
)

func Click() string {
	robotgo.Click()
	return fmt.Sprint("click")
}

func Move(x, y int) string {
	robotgo.Move(x, y)
	return fmt.Sprintf("Move x:%d, y:%d\n", x, y)
}

func MoveSmoothRelative(x, y int, args ...any) string {
	robotgo.MoveSmoothRelative(x, y, args)
	return fmt.Sprintf("MoveSmoothRelative x:%d, y:%d, %+v\n", x, y, args)
}
