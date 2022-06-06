package main

import (
	"bufio"
	"fmt"
	"github.com/CarsonSlovoka/go-controller/v2/app"
	"github.com/CarsonSlovoka/go-controller/v2/pkg/dll"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/go-vgo/robotgo"
	"os"
)

func main() {
	if err := app.InitApp(); err != nil {
		panic(err)
	}

	chHookDone := make(chan error) // 等待hook初始化完畢
	s := app.Script
	go s.Start(chHookDone)
	<-chHookDone
	hookEventDevice := robotgo.EventStart()

	suspendDeviceInfo := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		suspendDeviceInfo <- true // default do not show
		for {
			fmt.Println("Enter CMD: ")
			scanner.Scan()
			switch scanner.Text() {
			case "info":
				response, _ := dll.User32dll.MessageBox(0, "Device Info", "是否要顯示裝置動作資訊?", w32.MB_YESNO|w32.MB_TOPMOST|w32.MB_ICONQUESTION)
				if response == w32.IDYES {
					suspendDeviceInfo <- false
				} else {
					suspendDeviceInfo <- true
				}
			case "quit":
				s.Stop(false)
				return
			case "quit -f":
				s.Stop(true)
				return
			case "help":
				_, _ = dll.User32dll.MessageBox(0, "💡Help", `command list:
info
quit
quit -f
help
`,
					w32.MB_YESNO|w32.MB_TOPMOST|w32.MB_ICONINFORMATION)
			}
		}
	}()

	go func() {
		for {
			suspend, isOpen := <-suspendDeviceInfo
			if !isOpen {
				return
			}
			if suspend {
				continue
			}
			for e := range hookEventDevice { // 此通道可以得到滑鼠、鍵盤,...等訊息 // robotgo的通道是共用的，所以如果一直訪問channel，那麼已經註冊的hook.Event也會被影響，因為對於每一個event也是在等待channel來傳送消息，而您也無法指派該消息該由誰接收，變成先搶到的贏。因此我們需要制定一個中斷此訪問的策略
				// if e.Rawcode == w32.VK_F3 {} // F3可以開啟或關閉此功能
				select {
				case suspend, isOpen = <-suspendDeviceInfo:
					if !isOpen {
						return
					}
				default:
				}
				if suspend { // 當我們接收到暫停的命令，就跳離此迴圈來中斷channel的訪問
					break
				}
				// time.Sleep(1 * time.Second) // 預設只用50ms。睡眠會導致嚴重延遲。 如果刻意降速，雖然訊息並少，但是該發的event還是按照50ms去動作，也就是您在一秒內如果有20個動作，那麼channel還是會記錄這20個，您阻塞在這邊，只是慢慢抽取出來，該累機的訊息量依然造成累積，若您提取的速度不夠快就會導致延遲越來越嚴重。(當下時間看到的訊息可能是很久之前的)
				fmt.Println(e)
			}
		}
	}()
	<-robotgo.EventProcess(hookEventDevice) // 這列才會讓所有有註冊的Hook產生作用 // 此通道會阻塞直到robotgo.EventEnd()被觸發
}
