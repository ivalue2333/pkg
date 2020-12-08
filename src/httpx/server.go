package httpx

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ivalue2333/pkg/src/logx"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Name() (name string)
	GoTask(ctx context.Context) (err error)
	GetKernel() (kernel *gin.Engine)
}

type server struct {
	kernel  *gin.Engine
	options Options
}

func NewServer(options Options, opts ...Option) Server {
	// set default
	WithDefault(&options)

	for _, opt := range opts {
		opt(&options)
	}
	return NewServerWithOptions(options)
}

func NewServerWithOptions(options Options) Server {
	kernel := gin.New()

	// user set middles
	kernel.Use(options.Middles...)

	s := &server{kernel: kernel, options: options}
	return s
}

func (s *server) Name() string {
	return s.options.Name
}

// not very good
func (s *server) GoTask(ctx context.Context) error {
	srv := &http.Server{Handler: s.kernel, Addr: s.options.Address}

	// server listenAndServe
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logx.Errorf(ctx, "httpServer ListenAndServe err:%v", err)
		}
	}()

	// handle signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	ss := <-ch
	_ = srv.Shutdown(context.Background())
	logx.Infof(ctx, "Shutdown gracefully")
	return errors.New(fmt.Sprintf("Got signal %v, exit.", ss))
}

func (s *server) GetKernel() (kernel *gin.Engine) {
	return s.kernel
}
