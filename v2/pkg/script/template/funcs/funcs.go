package funcs

import "github.com/CarsonSlovoka/go-pkg/v2/tpl/funcs"

func AllFuncs() map[string]any {
	return map[string]any{
		// Log
		"Log": Log,

		// Time
		"Sleep":   Sleep,
		"MsSleep": MsSleep,

		// Key
		"TypeStr": TypeStr,

		// Hook
		"EventHook":  EventHook,
		"RunJob":     RunJob,
		"RunJobByID": RunJobByID,
		"ExitApp":    ExitApp,

		// Mouse
		"Click":              Click,
		"Move":               Move,
		"MoveSmoothRelative": MoveSmoothRelative,

		// Utils
		"dict": funcs.Dict,
	}
}
