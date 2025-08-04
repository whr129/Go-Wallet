package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	util "github.com/whr129/go-wallet/pkg/util"
)

func AuthMiddleware(ctx *gin.Context) {
	userID, ok := util.GetXUserID(ctx)
	if !ok {
		log.Printf("Missing or invalid X-User-ID header %d", userID)
		ctx.Abort()
		return
	}

	email, ok := util.GetXEmail(ctx)
	if !ok {
		log.Printf("Missing X-Email header for user ID: %s", email)
		ctx.Abort()
		return
	}

	role, ok := util.GetXRole(ctx)
	if !ok {
		log.Printf("Missing X-Role header for user ID: %d", userID)
		ctx.Abort()
		return
	}

	ctx.Set(util.X_USER_ID, userID)
	ctx.Set(util.X_EMAIL, email)
	ctx.Set(util.X_ROLE, role)

	ctx.Next()
}
