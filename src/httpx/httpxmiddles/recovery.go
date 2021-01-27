package httpxmiddles

import (
	"github.com/ivalue2333/pkg/src/logx"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() Middle {
	return RecoveryWithLogger(logx.StandardLogger())
}

func RecoveryWithLogger(l logx.Logger) Middle {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Errorf(c.Request.Context(), "panic recovered: err = %v, stack = %s", err, debug.Stack())
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
