package gpool

import (
	"context"
	"github.com/ivalue2333/pkg/logx"
	"runtime/debug"
	"sync"
	"time"
)

const (
	defaultGNum          = 30
	defaultChannelBuffer = 2000
	defaultInputRetry    = 3
)

type (
	Option func(w *Worker)
	ProcessFn func(data interface{})
)

func WithName(name string) Option {
	return func(w *Worker) {
		w.name = name
	}
}

func WithGNum(gNum int) Option {
	return func(w *Worker) {
		w.gNum = gNum
	}
}

func WithChannelBuffer(channelBuffer int) Option {
	return func(w *Worker) {
		w.channelBuffer = channelBuffer
	}
}

func WithInputRetry(retry int) Option {
	return func(w *Worker) {
		w.inputRetry = retry
	}
}

type Worker struct {
	name          string
	gNum          int
	channelBuffer int
	inputRetry    int
	dataChannel   chan interface{}
	process       ProcessFn
}

func NewWorkerWithOptions(process ProcessFn, options Options, opts2 ...Option) *Worker {
	opts := make([]Option, 0)

	if options.Name != "" {
		opts = append(opts, WithName(options.Name))
	}

	if options.GNum != 0 {
		opts = append(opts, WithGNum(options.GNum))
	}

	if options.ChannelBuffer != 0 {
		opts = append(opts, WithChannelBuffer(options.ChannelBuffer))
	}

	if options.InputRetry != 0 {
		opts = append(opts, WithInputRetry(options.InputRetry))
	}

	opts = append(opts, opts2...)

	return NewWorker(process, opts..., )
}

func NewWorker(process ProcessFn, opts ...Option) *Worker {
	worker := &Worker{
		gNum:          defaultGNum,
		channelBuffer: defaultChannelBuffer,
		inputRetry:    defaultInputRetry,
	}
	for _, opt := range opts {
		opt(worker)
	}
	worker.dataChannel = make(chan interface{}, worker.channelBuffer)
	worker.process = func(data interface{}) {
		defer func() {
			if err := recover(); err != nil {
				logx.Errorf(context.Background(), "panic happened: err = %v, stack = %s", err, debug.Stack())
			}
		}()
		process(data)
	}
	return worker
}

func (worker *Worker) Input(ctx context.Context, msg interface{}) {
	count := 0
	for {
		select {
		case worker.dataChannel <- msg:
			logx.Debugf(ctx, "input msg: %+v", msg)
			return
		default:
			if count >= worker.inputRetry {
				logx.Errorf(ctx, "input msg failed: %+v, channel full", msg)
				return
			}
			count += 1
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (worker *Worker) close() {
	close(worker.dataChannel)
}

func (worker *Worker) Name() string {
	return worker.name
}

func (worker *Worker) Run(ctx context.Context) (err error) {
	defer worker.close()
	wg := sync.WaitGroup{}
	for i := 0; i < worker.gNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			worker.run(ctx, index)
		}(i)
	}
	wg.Wait()
	return nil
}

func (worker *Worker) run(ctx context.Context, index int) {
	closeCount := 0
	closeFlag := false
	for {
		if closeFlag {
			closeCount += 1
			if closeCount > 5 {
				return
			}
		}
		select {
		case msg := <-worker.dataChannel:
			worker.process(msg)
		case <-ctx.Done():
			logx.Infof(ctx, "name(%s) g(%d) ctx.Done", worker.name, index)
			closeFlag = true
			ctx = context.Background()
		default:
			// default sleep for cpu ideal
			time.Sleep(100 * time.Millisecond)
		}
	}
}
