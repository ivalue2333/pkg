package httpcx

import (
	"context"
	"github.com/ivalue2333/pkg/src/logx"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func CreateRetryTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	wrapDialContext := func(ctx context.Context, network, address string) (conn net.Conn, err error) {

		for i := 0; i < 3; i++ {
			conn, err = dialer.DialContext(ctx, network, address)

			if IsConnectRefuseError(err) {
				logx.Warningf(ctx, "connection refused, err:%v, %v", err, i)
				continue
			}

			return conn, err
		}

		return
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           wrapDialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
}

func IsConnectRefuseError(err error) bool {
	if err == nil {
		return false
	}

	opError, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	const dialString = "dial"
	const netString = "tcp"
	isOk := opError.Op == dialString &&
		strings.Contains(opError.Err.Error(), "connection refused") &&
		opError.Net == netString

	return isOk
}
