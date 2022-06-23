package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	"simplebank/db/util"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetExchangeRate(t *testing.T) {
	randomAmount := strconv.FormatInt(util.RandomInt(1, 1000), 10)

	testCases := []struct {
		name             string
		mockToCurrency   string
		mockFromCurrency string
		mockAmount       string
		checkResponse    func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:             "OK",
			mockToCurrency:   util.HKD,
			mockFromCurrency: util.USD,
			mockAmount:       randomAmount,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:             "Invalid Currency",
			mockToCurrency:   "ABC",
			mockFromCurrency: util.USD,
			mockAmount:       randomAmount,
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

			exchangeURL := fmt.Sprintf("/exchange/%s/%s/%s", testCase.mockToCurrency, testCase.mockFromCurrency, testCase.mockAmount)
			request, err := http.NewRequest(http.MethodGet, exchangeURL, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}

}
