package mongox

import "time"

type (
	options struct {
		timeout time.Duration
	}

	Option func(opts *options)
)
