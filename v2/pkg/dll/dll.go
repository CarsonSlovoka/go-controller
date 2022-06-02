package dll

import "github.com/CarsonSlovoka/go-pkg/v2/w32"

var User32dll *w32.User32DLL

func init() {
	User32dll = w32.NewUser32DLL([]w32.ProcName{
		w32.PNMessageBox,
	})
}
