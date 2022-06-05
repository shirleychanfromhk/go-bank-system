package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type exchangeRequest struct {
	ToCurrency   string `json:"to" binding:"required"`
	FromCurrency string `json:"from" binding:"required"`
	Amount       string `json:"amount" binding:"required"`
}

type apiResponse struct {
	Status bool    `json:"success"`
	Result float64 `json:"result"`
}

func (server *Server) getExchangeRate(ctx *gin.Context) {
	var exchangeReq exchangeRequest
	if err := ctx.BindJSON(&exchangeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	url := "https://api.apilayer.com/exchangerates_data/convert?to=" + exchangeReq.ToCurrency + "&from=" +
		exchangeReq.FromCurrency + "&amount=" + exchangeReq.Amount

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	//The api key I will share in the ticket
	req.Header.Set("apikey", "shown in the ticket BK5333")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	var responseObject apiResponse
	json.Unmarshal(body, &responseObject)
	/* TODO check the status is it success
	if !apiResponse.Status {
		ctx.Status(400)
	}*/
	ctx.IndentedJSON(http.StatusOK, responseObject.Result)
}
