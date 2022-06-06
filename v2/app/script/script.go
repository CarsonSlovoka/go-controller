package script

import (
	"github.com/CarsonSlovoka/go-controller/v2/pkg/script/template/funcs"
	"os"
	"text/template"
	"time"
)

type runAble struct{}

func (r *runAble) execute(text string, context any) {
	t, err := new(template.Template).
		Funcs(funcs.AllFuncs()).
		Parse(text)
	if err != nil {
		panic(err)
	}
	if err = t.Execute(os.Stdout, context); err != nil {
		panic(err)
	}
}

// Template 為您輸入的json文件資料
type Template struct {
	Title, Desc string
	Hooks       []Hook
	Jobs        []*Job `json:"Jobs"`
}

type Hook struct {
	runAble
	Name, Desc string
	Func       string
}

type Job struct {
	runAble
	Name, Desc        string
	id                uint
	Loop              loop   `json:"Loop"`
	Async             bool   // 表示是否為異步的工作，若為false表示會等待此工作結束才會再往下進行下一個工作
	Func              string // 如果您的Job很單純，只需要執行一個Cmd，設定此Func即可
	WaitSignalToStart bool
	Cmd               []Cmd
	status                        // 目前的狀態
	chStatus          chan status // 要是可以寫入以及讀取的channel
}

type Cmd struct {
	Name, Desc string
	Func       string
}

type loop struct {
	MaxRun   int           // -1表示永久運行, 否則運行達到該次數就終止工作
	Interval time.Duration // 隔多久時間(毫秒)執行該命令一次 // 如果此項數值為-1，表示執行完之後就會等待通知才會再執行一次
}

// ----

type status uint

const RunForever = -1

const (
	StatusRunning status = iota
	StatusPaused  status = iota
	StatusStopped status = iota
)

// Script 用來運行腳本
type Script struct {
	*Template
	notify []chan status // 用來通知每一個job是否要啟動、暫停、終止
}
