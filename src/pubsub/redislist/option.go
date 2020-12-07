package redislist

type SubOption func(*sub)

func WithSubOptionTimeout(timeout int) SubOption {
	return func(s *sub) {
		s.pullTimeout = timeout
	}
}
