package middleware

import (
	"fmt"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/transport/http/core"
	"net/http"
	"runtime/debug"
)

func NewRecoverer(log logger.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						panic(rvr)
					}

					stack := debug.Stack()
					err := fmt.Errorf("panic recovered: %v\n%s", rvr, string(stack))

					core.WritePanic(w, log, err)
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
