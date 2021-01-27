package httpxmiddles

import (
	"bytes"
	"fmt"
	"github.com/ivalue2333/pkg/src/jsonx"
	"github.com/ivalue2333/pkg/src/logx"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const maxBodyLen = 1024

func LoggingRequest() Middle {
	return LoggingRequestWithLogger(logx.StandardLogger())
}

func LoggingRequestWithLogger(l logx.Logger) Middle {
	return func(c *gin.Context) {
		data := map[string]interface{}{
			"method": c.Request.Method,
			"uri":    c.Request.URL.RequestURI(),
			"body":   RequestBody(c),
		}
		l.Infof(c.Request.Context(), "incoming http request: data:%s", jsonx.MarshalUnsafe(data))
	}
}

func LoggingResponse() Middle {
	return LoggingResponseWithLogger(logx.StandardLogger())
}

func LoggingResponseWithLogger(l logx.Logger) Middle {
	return func(c *gin.Context) {
		rw := &responseWriter{Body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = rw
		now := time.Now()

		c.Next()

		usedTime := time.Since(now)
		statusCode := c.Writer.Status()
		statusText := http.StatusText(statusCode)
		body := rw.Body.Bytes()
		cnt := len(body)
		if cnt > maxBodyLen {
			cnt = maxBodyLen
		}

		data := map[string]interface{}{
			"status": fmt.Sprintf("%d %s", statusCode, statusText),
			"body":   string(body[:cnt]),
			"cost":   usedTime.Milliseconds(),
		}

		l.Infof(c.Request.Context(), "outgoing http response, data:%s", jsonx.MarshalUnsafe(data))
	}
}

func RequestBody(c *gin.Context) string {
	if c.Request.Body == nil || c.Request.Body == http.NoBody {
		return ""
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	_ = c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	cnt := len(body)
	if cnt > maxBodyLen {
		cnt = maxBodyLen
	}
	return string(body[:cnt])
}

type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w responseWriter) Write(body []byte) (int, error) {
	w.Body.Write(body)
	return w.ResponseWriter.Write(body)
}
