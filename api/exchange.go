package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"simplebank/db/util"

	"github.com/gin-gonic/gin"
)

const (
	URL = "https://api.apilayer.com/exchangerates_data/convert?to="
)

type exchangeRequest struct {
	ToCurrency   string `uri:"to" binding:"required"`
	FromCurrency string `uri:"from" binding:"required"`
	Amount       string `uri:"amount" binding:"required"`
}

type apiResponse struct {
	Status bool    `json:"success"`
	Result float64 `json:"result"`
	Error  apiError
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (server *Server) getExchangeRate(ctx *gin.Context) {
	var exchangeReq exchangeRequest
	if err := ctx.ShouldBindUri(&exchangeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	url := buildApiUrl(exchangeReq)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	config, err := util.LoadViberConfig("../")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	req.Header.Set("apikey", config.ApiKey)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	var responseObject apiResponse
	json.Unmarshal(body, &responseObject)
	if err != nil {
		return
	}
	if res.StatusCode == 400 {
		ctx.IndentedJSON(http.StatusBadRequest, responseObject.Error.Message)
	}

	if res.StatusCode == 429 {
		ctx.IndentedJSON(http.StatusTooManyRequests, responseObject.Error.Message)
	}

	if res.StatusCode != 200 {
		ctx.IndentedJSON(http.StatusInternalServerError, "Internal server error")
	}

	log.Printf(string(body))
	ctx.IndentedJSON(http.StatusOK, responseObject.Result)
}

func buildApiUrl(exchangeReq exchangeRequest) string {
	result := URL + exchangeReq.ToCurrency + "&from=" +
		exchangeReq.FromCurrency + "&amount=" + exchangeReq.Amount

	return result
}
