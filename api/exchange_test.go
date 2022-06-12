package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	"simplebank/db/util"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetExchangeRate(t *testing.T) {
	mockToCurrency := util.HKD
	mockFromCurrency := util.USD
	mockAmount := strconv.FormatInt(util.RandomInt(1, 1000), 10)

	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"to":     mockToCurrency,
				"from":   mockFromCurrency,
				"amount": mockAmount,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid Currency",
			body: gin.H{
				"to":     "ABC",
				"from":   mockFromCurrency,
				"amount": mockAmount,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			res, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			exchangeURL := "/exchange"
			request, err := http.NewRequest(http.MethodGet, exchangeURL, bytes.NewReader(res))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}

}
