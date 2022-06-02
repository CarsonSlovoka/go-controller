package funcs

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"reflect"
)

type Job interface {
	CheckName(name string) error
	CheckID(id uint) error
	Run()
	Wait()
	Stop()
	IsClosed() bool
}

func EventHook(values ...any) string {
	if len(values) < 3 {
		panic("insufficient parameters (length should be greater than 2)")
	}
	var (
		inputKeyPressed []any
		inputWhen       string
	)

	inputWhen = values[0].(string)              // 0
	inputKeyPressed = values[1 : len(values)-1] // 1 ~ -1
	callback := values[len(values)-1].(func())  // -1
	when, exists := map[string]uint8{
		"KeyDown": hook.KeyDown,
		"KeyHold": hook.KeyHold,
		"KeyUp":   hook.KeyUp,

		"MouseUp":   hook.MouseUp,
		"MouseHold": hook.MouseHold,
		"MouseDown": hook.MouseDown,
		"MouseMove": hook.MouseMove,
		"MouseDrag": hook.MouseDrag,
	}[inputWhen]

	if !exists {
		panic("'when' is not exist.")
	}
	keysPressed := make([]string, len(inputKeyPressed))
	for i, keyName := range inputKeyPressed {
		switch reflect.ValueOf(keyName).Kind() {
		case reflect.String:
			keysPressed[i] = keyName.(string)
		default:
			panic(`type must equal to "string"`)
		}
	}

	robotgo.EventHook(when, keysPressed, func(e hook.Event) { // register hook event // 如果您用多個routine來註冊，要確保一個接一個，不能註冊中途有人插隊，不然會發生map conflict
		callback()
	})
	return ""
}

func filterJob[T uint | string](jobs []Job, key T) *Job {
	for _, job := range jobs {
		if fmt.Sprintf("%T", key) == "string" {
			if job.CheckName(any(key).(string)) == nil {
				return &job
			}
		} else {
			if job.CheckID(any(key).(uint)) == nil {
				return &job
			}
		}
	}
	return nil
}

func RunJobByID(jobID uint, jobs []Job) func() {
	return func() {
		job := filterJob(jobs, jobID)
		if job == nil {
			panic(fmt.Sprintf("jobID:%d not found", jobID))
		}
		(*job).Run()
	}
}

func RunJob(name string, jobs []Job) func() {
	return func() {
		job := filterJob(jobs, name)
		fmt.Printf("%+v", job)
		if job == nil {
			panic(fmt.Sprintf("job name:%q not found", name))
		}
		(*job).Run()
	}
}

func ExitApp(jobs []Job) func() { // 通知所有工作都結束作業就能終止程式
	return func() {
		for _, job := range jobs {
			if job.IsClosed() {
				continue
			}
			job.Stop()
		}
		robotgo.EventEnd()
	}
}
