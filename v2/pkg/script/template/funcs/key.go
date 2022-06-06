package funcs

import (
	"fmt"
	"github.com/go-vgo/robotgo"
)

func TypeStr(str string, args ...float64) string {
	robotgo.TypeStr(str, args...)
	if len(args) > 0 {
		return fmt.Sprintf("TypeStr %s %v\n", str, args)
	}
	return fmt.Sprintf("TypeStr %s", str)
}
