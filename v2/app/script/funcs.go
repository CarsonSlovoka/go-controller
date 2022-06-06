package script

import (
	"encoding/json"
	"fmt"
	"github.com/CarsonSlovoka/go-controller/v2/pkg/script/template/funcs"
	"github.com/go-vgo/robotgo"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func (s *status) String() string {
	switch *s {
	case StatusRunning:
		return "Running"
	case StatusPaused:
		return "Paused"
	case StatusStopped:
		return "Stopped"
	}
	return ""
}

// ----

func NewTemplate(filePath string) (*Template, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	dataBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	t := new(Template)
	err = json.Unmarshal(dataBytes, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func NewScript(t *Template) *Script {
	s := Script{Template: t}
	s.notify = make([]chan status, len(t.Jobs))
	return &s
}

func (s *Script) Start(out chan<- error) {
	contextJobs := make([]funcs.Job, len(s.Jobs))
	for i, job := range s.Jobs {
		contextJobs[i] = job
	}

	// Hook
	{
		mux := sync.Mutex{} // 註冊hook.Event註冊發生衝突(map conflict)
		for _, hook := range s.Hooks {
			go func(h Hook) {
				mux.Lock()
				defer mux.Unlock()
				h.Execute(contextJobs)
				log.Printf("init hook %q done\n", h.Name)
			}(hook)
		}
		out <- nil // 通知外層所有hook已經初始化完畢
	}

	// Job
	{
		wg := sync.WaitGroup{}
		for i, job := range s.Jobs {
			job.id = uint(i) // init job ID
			s.notify[i] = make(chan status)
			job.status = StatusPaused // 開始都是等待的狀態
			job.chStatus = s.notify[i]
			if !job.Async { // 若不為異步作業，就會等待此工作結束才會執行下一個工作
				wg.Add(1)
			}
			go func(j *Job) {
				// job.Execute(nil) // ❗ 注意在for中做併發時，如果直接用for的變數會有風險，當變數改變時，routine內使用到的也會跟著變，這可能不是您所期望的，因此最好把它當作參數傳進來
				j.Execute(nil) // Cmd的工作暫時不傳入任何的context
				if !job.Async {
					wg.Done()
				}
				log.Printf("job.Execute done. [%q]\n", job.Name)
			}(job)
			if !job.WaitSignalToStart {
				job.chStatus <- StatusRunning
			}
			if !job.Async {
				wg.Wait()
			}
		}
	}

	// out <- nil
}

// ----

func (h *Hook) Execute(context any) {
	if text := h.Func; text != "" {
		h.execute(text, context)
	}
}

// ----

func (job *Job) CheckName(name string) error {
	if job.Name != name {
		return fmt.Errorf("it is different")
	}
	return nil
}

func (job *Job) CheckID(id uint) error {
	if job.id != id {
		return fmt.Errorf("it is different")
	}
	return nil
}

func (job *Job) Wait() {
	job.status = StatusPaused
	log.Printf("%q [%v] waiting...\n", job.Name, job.chStatus)
	s, isOpen := <-job.chStatus // 等待通知
	oldStatus := job.status
	if !isOpen {
		job.status = StatusStopped
	} else {
		job.status = s
	}
	if oldStatus != job.status {
		log.Printf("[%q] The status has changed from %q to %q.\n", job.Name, oldStatus.String(), job.status.String())
	}
}

func (job *Job) Run() {
	if job.status == StatusStopped {
		log.Printf("%s Already closed\n", job.Name)
		return
	}
	job.chStatus <- StatusRunning
}

func (job *Job) Pause() {
	if job.status == StatusRunning {
		job.status = StatusPaused
		log.Printf("pause job: %q\n", job.Name)
		return
	}
}

func (job *Job) Stop() {
	if job.status == StatusStopped {
		log.Printf("%s Already closed\n", job.Name)
		return
	}
	log.Printf("%s Stop\n", job.Name)
	job.status = StatusStopped
	close(job.chStatus)
}

func (job *Job) IsClosed() bool {
	if job.status == StatusStopped {
		return true
	}
	return false
}

func (job *Job) Execute(context any) {
	for {
		job.Wait()
		if job.status == StatusRunning {
			break
		}
		if job.status == StatusStopped {
			return
		}
	}

	ExecuteFunc := func(ctx any) status {
		if text := job.Func; text != "" { // 對於有指定Func的Job，視為簡單的工作，直接運行，忽略所有Cmd的項目
			job.execute(text, ctx)
		} else if len(job.Cmd) > 0 {
			for idx, curCmd := range job.Cmd {
				if job.status == StatusPaused { // 如果cmd太多，可以有辦法中途暫停
					if nextIdx := idx + 1; nextIdx < len(job.Cmd) {
						nextCmd := job.Cmd[nextIdx]
						var nextCmdDesc string
						if desc := nextCmd.Name; desc != "" {
							nextCmdDesc = desc
						} else if desc = nextCmd.Desc; desc != "" {
							nextCmdDesc = desc
						} else {
							nextCmdDesc = nextCmd.Func
						}
						log.Printf("The %q has paused. The next command is %q", job.Name, nextCmdDesc)
					}
					job.Wait()
				}
				if job.status == StatusStopped {
					return StatusStopped
				}
				job.execute(curCmd.Func, ctx)
			}
		} else {
			log.Println("This job neither contains 'Func' nor 'Cmd'(or empty), so it is an empty job.")
			return StatusStopped
		}

		if job.Loop.Interval == -1 { // 等待通知才會再執行一次
			job.Wait()
			return job.status // 返回通知後的狀態
		} else if job.Loop.MaxRun == 0 || job.Loop.MaxRun == 1 { // 如果是只有執行一次或者預設值(0)，就直接終止
			return StatusStopped
		} else {
			time.Sleep(job.Loop.Interval * time.Millisecond)
			return StatusRunning
		}
	}

	maxRun := job.Loop.MaxRun
	count := 0
	for {
		count++
		s := ExecuteFunc(context)
		if s == StatusStopped ||
			(count >= maxRun && maxRun != RunForever) {
			break
		}
		log.Printf("Has run the job %q %d times.", job.Name, count)
	}
	job.Stop()
	return
}

func (s *Script) Stop(force bool) {
	if !force {
		for _, job := range s.Template.Jobs {
			job.Stop()
		}
	}
	robotgo.EventEnd()
}
