package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/whr129/go-wallet/cmd/wallet-service/db/sqlc"
	"github.com/whr129/go-wallet/cmd/wallet-service/dto"
	util "github.com/whr129/go-wallet/pkg/util"
)

func (server *Server) createAccount(ctx *gin.Context) {
	var req dto.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authUserID := ctx.MustGet(util.X_USER_ID)
	userID, ok := authUserID.(int64)
	log.Printf("Creating account for user ID: %d", userID)
	if !ok {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user ID type")))
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user ID")))
		return
	}

	id, err := util.GenerateID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to generate account ID)")))
		return
	}

	arg := db.CreateAccountParams{
		ID:       id,
		UserID:   userID,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Printf("Account created successfully: %+v", account)
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req dto.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authUserID := ctx.MustGet(util.X_USER_ID)
	if account.UserID != authUserID.(int64) {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req dto.ListAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authUserID := ctx.MustGet(util.X_USER_ID)
	arg := db.ListAccountsParams{
		UserID: authUserID.(int64),
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
