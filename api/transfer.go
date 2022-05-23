package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
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

	validFromAccount, fromAccount := server.validAccount(ctx, req.FromAccountID)
	if !validFromAccount || !server.validCurrency(ctx, fromAccount, req.Currency) {
		return
	}

	validToAccount, toAccount := server.validAccount(ctx, req.ToAccountID)
	if !validToAccount || !server.validCurrency(ctx, toAccount, req.Currency) {
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

func (server *Server) validAccount(ctx *gin.Context, accountID int64) (bool, db.Account) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false, db.Account{}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false, db.Account{}
	}
	return true, account
}

func (server *Server) validCurrency(ctx *gin.Context, account db.Account, currency string) bool {
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
