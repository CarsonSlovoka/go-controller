package funcs

import (
	"fmt"
	"github.com/go-vgo/robotgo"
)

// Click (button, double)
//	Examples:
//		Click("left", true) // double click left button
//		Click("right", false)
func Click(args ...string) string {
	button := "left"
	isDouble := false
	switch len(args) {
	case 0:
		robotgo.Click()
	default:
		fallthrough
	case 2:
		if args[1] == "true" || args[1] == "1" {
			isDouble = true
		}
		fallthrough
	case 1:
		button = args[0]
		robotgo.Click(button, isDouble)
	}
	return fmt.Sprintf("click %s double:%v\n", button, isDouble)
}

func Move(x, y int) string {
	robotgo.Move(x, y)
	return fmt.Sprintf("Move x:%d, y:%d\n", x, y)
}

// MoveSmooth (x, y int, low, high float64, mouseDelay int)
//  Examples:
//		MoveSmooth(10, 10)
//		MoveSmooth(10, 10, 1.0, 2.0)
func MoveSmooth(x, y int, args ...any) string {
	robotgo.MoveSmooth(x, y, args)
	return fmt.Sprintf("MoveSmooth x:%d, y:%d, %+v\n", x, y, args)
}

func MoveSmoothRelative(x, y int, args ...any) string {
	robotgo.MoveSmoothRelative(x, y, args)
	return fmt.Sprintf("MoveSmoothRelative x:%d, y:%d, %+v\n", x, y, args)
}
