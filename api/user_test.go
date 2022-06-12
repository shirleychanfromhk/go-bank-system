package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/db/util"
	"testing"
)

func userMatcher(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FirstName, gotUser.FirstName)
	require.Equal(t, user.LastName, gotUser.LastName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.ContactNumber, gotUser.ContactNumber)
	require.Equal(t, user.Address, gotUser.Address)
	require.Empty(t, gotUser.HashedPassword)
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		FirstName:      util.RandomUsername(),
		LastName:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		ContactNumber:  util.RandomContactNumber(),
		Address:        util.RandomAddress(),
	}
	return
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.ValidPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":       user.Username,
				"password":       password,
				"first_name":     user.FirstName,
				"last_name":      user.LastName,
				"email":          user.Email,
				"contact_number": user.ContactNumber,
				"address":        user.Address,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:      user.Username,
					FirstName:     user.FirstName,
					LastName:      user.LastName,
					Email:         user.Email,
					ContactNumber: user.ContactNumber,
					Address:       user.Address,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				userMatcher(t, recorder.Body, user)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":       user.Username,
				"password":       password,
				"first_name":     user.FirstName,
				"last_name":      user.LastName,
				"email":          user.Email,
				"contact_number": user.ContactNumber,
				"address":        user.Address,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":       "!Invalidusername#",
				"password":       password,
				"first_name":     user.FirstName,
				"last_name":      user.LastName,
				"email":          user.Email,
				"contact_number": user.ContactNumber,
				"address":        user.Address,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":       user.Username,
				"password":       password,
				"first_name":     user.FirstName,
				"last_name":      user.LastName,
				"email":          "invalidemail",
				"contact_number": user.ContactNumber,
				"address":        user.Address,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":       user.Username,
				"password":       "123",
				"first_name":     user.FirstName,
				"last_name":      user.LastName,
				"email":          user.Email,
				"contact_number": user.ContactNumber,
				"address":        user.Address,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
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
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			res, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			userURL := "/users"
			request, err := http.NewRequest(http.MethodPost, userURL, bytes.NewReader(res))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}

}
