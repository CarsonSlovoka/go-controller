package script

import (
	"encoding/json"
	"fmt"
	"github.com/CarsonSlovoka/go-controller/v2/pkg/script/template/funcs"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

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
			if job.Async { // 若為異步作業，就會等待此工作結束才會執行下一個工作
				wg.Add(1)
			}
			go func() {
				job.Execute(contextJobs)
				if job.Async {
					wg.Done()
				}
			}()
			if !job.WaitSignalToStart {
				job.chStatus <- StatusRunning
			}
			if job.Async {
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
	/*
		if job.status == StatusStopped { // 由於job有可能由其他routine重複執行，所以要確認
			log.Printf("%q Already closed", job.Name)
			return
		}
		if job.status == StatusPaused {
			log.Printf("%q Already waiting", job.Name)
			return
		}
	*/
	job.status = StatusPaused
	s, isOpen := <-job.chStatus // 等待通知
	if !isOpen {
		job.status = StatusStopped
		return
	}
	job.status = s // 將status更新完通知的狀態
}

func (job *Job) Run() {
	/*
		if job.status == StatusStopped {
			log.Printf("%q You can't run again when the job has been closed.", job.Name)
			return
		}
		if job.status == StatusRunning {
			log.Printf("%q Already running", job.Name)
			return
		}
	*/
	job.chStatus <- StatusRunning
}

func (job *Job) Stop() {
	if job.status == StatusStopped {
		log.Printf("%s Already closed", job.Name)
		return
	}

	job.status = StatusStopped
	close(job.chStatus)
}

func (job *Job) IsClosed() bool {
	if job.status == StatusStopped {
		return true
	}
	return false
}

/*
func (job *Job) Response(s status) {
	job.status = s
	job.chStatus <- s
}
*/

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

	if text := job.Func; text != "" { // 對於有指定Func的Job，視為簡單的工作，直接運行後就終了
		job.execute(text, context)
		job.Stop()
		return
	}

	if len(job.Cmd) == 0 {
		log.Println("This job neither contains 'Func' nor 'Cmd'(or empty), so it is an empty job.")
		job.Stop()
		return
	}

	ExecuteCMD := func(ctx any) status {
		if job.status == StatusPaused {
			job.Wait()
		}
		if job.status == StatusStopped {
			return StatusStopped
		}
		for _, curCmd := range job.Cmd {
			job.execute(curCmd.Func, ctx)
		}
		time.Sleep(job.Loop.Interval)
		return StatusRunning
	}

	maxRun := job.Loop.MaxRun
	if maxRun == RunForever {
		for {
			s := ExecuteCMD(nil)
			if s == StatusStopped {
				job.Stop()
				return
			}
		}
	}
	for count := 0; count < maxRun; count++ {
		s := ExecuteCMD(nil)
		if s == StatusStopped {
			job.Stop()
			return
		}
	}
	return
}
