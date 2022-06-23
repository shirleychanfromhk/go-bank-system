package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransaction(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	validFromAccount, fromAccount := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !validFromAccount {
		return
	}

	if ctx.MustGet(authPayLoadKey).(*token.Payload).Username != fromAccount.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("From account does not belong to the user")))
		return
	}

	validToAccount, _ := server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !validToAccount {
		return
	}

	arg := db.TransactionTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transaction, err := server.store.TransactionTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (bool, db.Account) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false, account
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false, account
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false, account
	}
	return true, account
}
