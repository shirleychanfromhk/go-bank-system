package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/db/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	checkUser := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		FirstName:      util.RandomUsername(),
		LastName:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		ContactNumber:  util.RandomContactNumber(),
		Address:        util.RandomAddress(),
	}

	testUser, err := testQueries.CreateUser(context.Background(), checkUser)

	require.NoError(t, err)
	require.NotEmpty(t, testUser)

	require.Equal(t, checkUser.Username, testUser.Username)
	require.Equal(t, checkUser.HashedPassword, testUser.HashedPassword)
	require.Equal(t, checkUser.FirstName, testUser.FirstName)
	require.Equal(t, checkUser.LastName, testUser.LastName)
	require.Equal(t, checkUser.Email, testUser.Email)
	require.Equal(t, checkUser.ContactNumber, testUser.ContactNumber)
	require.Equal(t, checkUser.Address, testUser.ContactNumber)

	require.NotZero(t, testUser.CreatedAt)
	require.NotZero(t, testUser.UpdatedAt)

	return testUser
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.ContactNumber, user2.ContactNumber)
	require.Equal(t, user1.Address, user2.ContactNumber)

	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
