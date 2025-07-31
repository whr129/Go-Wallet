package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/whr129/go-wallet/pkg/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	X_USER_ID               = "X-User-ID"
	X_EMAIL                 = "X-Email"
	X_ROLE                  = "X-Role"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker, redisClient redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		authSessionDetails, err := redisClient.Get(ctx, accessToken).Result()

		if err == redis.Nil {
			err = fmt.Errorf("access token %s not found", accessToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		} else if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		userInfo := []byte(authSessionDetails)

		var payload AuthSessionDetails

		err = json.Unmarshal(userInfo, &payload)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(X_USER_ID, payload.UserID)
		ctx.Set(X_EMAIL, payload.Email)
		ctx.Set(X_ROLE, payload.Role)
		ctx.Next()
	}
}
