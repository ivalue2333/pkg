package runner

import (
	"context"
	"github.com/ivalue2333/pkg/logx"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Runner interface {
	RunTasks(ctx context.Context) (err error)
}

type runner struct {
	tasks []Task
}

func NewRunner(tasks ...Task) Runner {
	return &runner{tasks: tasks}
}

func (r *runner) RunTasks(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(quit)
		close(quit)
	}()
	go func() {
		<-quit;
		cancel()
	}()

	var wg sync.WaitGroup
	for _, t := range r.tasks {
		task := t
		wg.Add(1)
		go func() {
			defer wg.Done()
			logx.Infof(ctx, "task(%s) is starting", task.Name())
			if err := task.GoTask(ctx); err != nil {
				logx.Errorf(ctx, "task(%s) run with error(%v)", task.Name(), err)
			}
			logx.Infof(ctx, "task(%s) is stopped", task.Name())
			cancel()
		}()
	}
	wg.Wait()

	return nil
}
