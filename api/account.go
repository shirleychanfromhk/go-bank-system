package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type updateAccountRequest struct {
	Id       int64  `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Balance  int64  `json:"balance"  binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Location string `json:"location" binding:"required"`
}

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
	Location string `json:"location" binding:"required"`
}

type accountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}

// Define a GraphQL schema using your GraphQL library
type Account struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
}

// var rootQuery = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "Query",
// 	Fields: graphql.Fields{
// 		"getAccount": &graphql.Field{
// 			Type: accountType,
// 			Args: graphql.FieldConfigArgument{
// 				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
// 			},
// 			Resolve: getAccountResolver, // You need to define this function
// 		},
// 	},
// })

// GraphQL resolver for Query type
type Resolver struct {
	store db.Store // Assuming db.Store is your data store interface
}

func (r *Resolver) GetAccount(ctx *gin.Context, args struct{ ID int64 }) (*Account, error) {
	account, err := r.store.GetAccount(ctx, args.ID)
	ginCtx := ctx.Value("ginContext").(*gin.Context)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Account not found")

		}
		return nil, err
	}

	if account.Username != ginCtx.MustGet(authPayLoadKey).(*token.Payload).Username {
		// Handle unauthorized access error
		return nil, errors.New("The account does not belong to the user")
	}

	// Convert your data model to GraphQL type
	graphqlAccount := &Account{
		ID:        account.ID,
		Username:  account.Username,
		Balance:   account.Balance,
		Currency:  account.Currency,
		Location:  account.Location,
		CreatedAt: account.CreatedAt.Format("2024-01-02T15:04:05Z"),
	}

	return graphqlAccount, nil
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req accountByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Username != ctx.MustGet(authPayLoadKey).(*token.Payload).Username {
		err := errors.New("The account does not belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getListAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Username: ctx.MustGet(authPayLoadKey).(*token.Payload).Username,
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Username: ctx.MustGet(authPayLoadKey).(*token.Payload).Username,
		Currency: req.Currency,
		Location: req.Location,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:       req.Id,
		Username: req.Username,
		Balance:  req.Balance,
		Currency: req.Currency,
		Location: req.Location,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req accountByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Delete Successed.")
}
