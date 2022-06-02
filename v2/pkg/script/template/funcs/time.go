package funcs

import (
	"fmt"
	"time"
)

func Sleep(sec time.Duration) string {
	time.Sleep(sec * time.Second)
	return fmt.Sprintf("sleep: %d sec", sec)
}

func MsSleep(ms time.Duration) string {
	time.Sleep(ms * time.Millisecond)
	return fmt.Sprintf("sleep: %d ms", ms)
}
