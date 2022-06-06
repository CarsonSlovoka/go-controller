package funcs

import (
	"github.com/CarsonSlovoka/go-controller/v2/pkg/dll"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
)

func MessageBoxOK(caption, content string) string {
	_, _ = dll.User32dll.MessageBox(0, caption, content, w32.MB_OK|w32.MB_ICONINFORMATION|w32.MB_TOPMOST)
	return ""
}
