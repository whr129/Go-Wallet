package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	X_USER_ID = "X-User-ID"
	X_EMAIL   = "X-Email"
)

func GetXUserID(ctx *gin.Context) (int64, bool) {
	xUserID := ctx.GetHeader(X_USER_ID)
	userID, err := strconv.ParseInt(xUserID, 10, 64)
	if err != nil {
		return 0, false
	}
	return userID, true
}

func GetXEmail(ctx *gin.Context) (string, bool) {
	xEmail := ctx.GetHeader(X_EMAIL)
	if xEmail == "" {
		return "", false
	}
	return xEmail, true
}
