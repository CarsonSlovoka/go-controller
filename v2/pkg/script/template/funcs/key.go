package funcs

import (
	"fmt"
	"github.com/go-vgo/robotgo"
)

// TypeStr send a string, support UTF-8
// USAGE:
// 	TypeStr(string: The string to send, float64: microSleep time, x11 option)
// Examples:
// 	TypeStr("Hello World!")
//	TypeStr("Hello World!", 1.0)
func TypeStr(str string, args ...float64) string {
	robotgo.TypeStr(str, args...)
	if len(args) > 0 {
		return fmt.Sprintf("TypeStr %s %v\n", str, args)
	}
	return fmt.Sprintf("TypeStr %s", str)
}

// KeyTap tap the keyboard code
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
//  https://github.com/CarsonSlovoka/robotgo/blob/v0.100.10/docs/keys.md
// Examples:
//	KeyTap("a", "control")
//  KeyTap("c", "control")
//  KeyTap("end")
//	KeyTap("enter")
//  KeyTap("v", "control")
//  ----
//	KeyTap("f1", "control", "alt")
//	KeyTap("i", []string{"alt", "command"})
func KeyTap(tapKey string, args ...string) string {
	if len(args) > 0 {
		robotgo.KeyTap(tapKey, args)
	} else {
		robotgo.KeyTap(tapKey)
	}
	return fmt.Sprintf("KeyTap %s %v\n", tapKey, args)
}
