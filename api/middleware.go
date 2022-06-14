package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"simplebank/token"
	"strings"
)

const authTypeBearer = "bearer"
const authHeaderKey = "authorization"
const authPayLoadKey = "authorization_payload"

const noAuthHeader = "Authorization header is not provided."
const invalidAuthHeader = "Invalid authorization header format."
const unSupportedAuth = "Server does not support the current auth method"

func errRes(errorMessage string) gin.H {
	err := errors.New(errorMessage)
	return errorResponse(err)
}

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authHeaderKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(noAuthHeader))
			return
		}

		tokenFields := strings.Fields(authorizationHeader)
		if len(tokenFields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(invalidAuthHeader))
		}

		authType := strings.ToLower(tokenFields[0])
		if authType != authTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(unSupportedAuth))
		}

		accessToken := tokenFields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(unSupportedAuth))
		}

		ctx.Set(authPayLoadKey, payload)
		ctx.Next()
	}
}
