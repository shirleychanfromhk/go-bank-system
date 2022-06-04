package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestCorrectPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = ValidPassword(password, hashedPassword)
	require.NoError(t, err)
}

func TestWrongPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = ValidPassword(RandomString(6), hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

func TestDoubleHashing(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = ValidPassword(password, hashedPassword1)
	require.NoError(t, err)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
