package middleware

import (
	"context"
	"github.com/thxhix/chat/internal/domain/token"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/security/jwt"
	"github.com/thxhix/chat/internal/transport/http/core"

	"net/http"
	"strconv"
	"strings"
)

type TokenParser interface {
	ParseAccessToken(tokenStr string) (userID string, err error)
}

type ctxKey string

const CtxKeyUserID ctxKey = "user_id"

func Authorize(logger logger.ILogger, jwtManager jwt.IJWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				core.WriteError(w, logger, token.ErrMissingToken)
				return
			}

			providedToken := strings.TrimPrefix(h, "Bearer ")

			userID, err := jwtManager.ParseAccessToken(providedToken)
			if err != nil {
				core.WriteError(w, logger, token.ErrDeadToken)
				return
			}

			ctx := context.WithValue(r.Context(), CtxKeyUserID, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromCtx(ctx context.Context) (int64, bool) {
	v := ctx.Value(CtxKeyUserID)
	str, ok := v.(string)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
