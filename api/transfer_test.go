package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/db/util"
	"testing"
)

func TestTransactionAPI(t *testing.T) {
	mockAmount := int64(10)
	mockCurrency := util.HKD

	mockAccount1 := randomAccount(t)
	mockAccount2 := randomAccount(t)
	mockAccount1.Currency = mockCurrency
	mockAccount2.Currency = mockCurrency

	mockCurrencyMisMatchAccount := randomAccount(t)
	mockCurrencyMisMatchAccount.Currency = util.USD

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": mockAccount1.ID,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        util.HKD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount1.ID)).Times(1).Return(mockAccount1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount2.ID)).Times(1).Return(mockAccount2, nil)

				arg := db.TransactionTxParams{
					FromAccountID: mockAccount1.ID,
					ToAccountID:   mockAccount2.ID,
					Amount:        mockAmount,
				}
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "From Account Currency MisMatch",
			body: gin.H{
				"from_account_id": mockCurrencyMisMatchAccount.ID,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        util.HKD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockCurrencyMisMatchAccount.ID)).Times(1).Return(mockCurrencyMisMatchAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount2.ID)).Times(0).Return(mockAccount2, nil)

				arg := db.TransactionTxParams{
					FromAccountID: mockCurrencyMisMatchAccount.ID,
					ToAccountID:   mockAccount2.ID,
					Amount:        mockAmount,
				}
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "From Account Currency MisMatch",
			body: gin.H{
				"from_account_id": mockAccount1.ID,
				"to_account_id":   mockCurrencyMisMatchAccount.ID,
				"amount":          mockAmount,
				"currency":        util.HKD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount1.ID)).Times(1).Return(mockAccount1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockCurrencyMisMatchAccount.ID)).Times(1).Return(mockCurrencyMisMatchAccount, nil)

				arg := db.TransactionTxParams{
					FromAccountID: mockAccount1.ID,
					ToAccountID:   mockCurrencyMisMatchAccount.ID,
					Amount:        mockAmount,
				}
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Currency",
			body: gin.H{
				"from_account_id": mockAccount1.ID,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        "BTC",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount1.ID)).Times(0).Return(mockAccount1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount2.ID)).Times(0).Return(mockAccount2, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccount not found",
			body: gin.H{
				"from_account_id": mockAccount1.ID,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        mockCurrency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount2.ID)).Times(0)
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccount not found",
			body: gin.H{
				"from_account_id": mockAccount1.ID,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        mockCurrency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount1.ID)).Times(1).Return(mockAccount1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Account Error",
			body: gin.H{
				"from_account_id": mockAccount1.ID + 1,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        mockCurrency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: gin.H{
				"from_account_id": mockAccount1.ID,
				"to_account_id":   mockAccount2.ID,
				"amount":          mockAmount,
				"currency":        mockCurrency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount1.ID)).Times(1).Return(mockAccount1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(mockAccount2.ID)).Times(1).Return(mockAccount2, nil)
				store.EXPECT().TransactionTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransactionTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			res, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			transactionURL := "/transactions"
			request, err := http.NewRequest(http.MethodPost, transactionURL, bytes.NewReader(res))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
