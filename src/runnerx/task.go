package runnerx

import "context"

type Task interface {
	Name() string
	GoTask(ctx context.Context) (err error)
}
