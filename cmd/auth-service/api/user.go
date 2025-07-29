package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	db "github.com/whr129/go-wallet/cmd/auth-service/db/sqlc"
	"github.com/whr129/go-wallet/cmd/auth-service/dto"
	token "github.com/whr129/go-wallet/pkg/token"
	util "github.com/whr129/go-wallet/pkg/util"
)

func newUserResponse(user db.User) dto.UserResponse {
	return dto.UserResponse{
		ID:                user.ID,
		UserName:          user.UserName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Create user_id
	id, err := util.GenerateID()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		ID:           id,
		UserName:     req.UserName,
		HashPassword: hashedPassword,
		Email:        req.Email,
		Role:         req.Role,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if the user is in redis
	userCheckMsg := fmt.Sprintf("user:%d", user.ID)
	_, err = server.redisClient.Get(ctx, userCheckMsg).Result()

	if err != nil && err != redis.Nil {
		// User is already logged in, return an error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else if err == nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User Already Logged In"})
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.UserName,
		user.Email,
		user.Role,
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.UserName,
		user.Email,
		user.Role,
		server.config.RefreshTokenDuration,
		token.TokenTypeRefreshToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, err := util.GenerateID()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           id,
		UserID:       refreshPayload.UserID,
		UserName:     user.UserName,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Store the session in Redis
	userMsg := fmt.Sprintf("user:%d", user.ID)
	authSession := dto.AuthSessionDetails{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: session.ID,
		Role:      user.Role,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	var authSessionDetails []byte
	authSessionDetails, err = json.Marshal(authSession)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.redisClient.Set(ctx, userMsg, authSessionDetails, server.config.RefreshTokenDuration).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Store the refresh token in Redis
	refreshMsg := fmt.Sprintf("refresh:%s", refreshToken)
	err = server.redisClient.Set(ctx, refreshMsg, authSessionDetails, server.config.RefreshTokenDuration).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Store the access token in Redis
	accessMsg := fmt.Sprintf("access:%s", accessToken)
	err = server.redisClient.Set(ctx, accessMsg, authSessionDetails, server.config.AccessTokenDuration).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := dto.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
