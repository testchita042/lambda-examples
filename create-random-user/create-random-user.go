package createrandomuser

import (
	"fmt"
	"math/rand"
)

type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

func CreateRandomUser() (User, error) {

	randomN := 1e6 + rand.Intn(900000)

	return User{
		Username:    fmt.Sprintf("user-%d", randomN),
		DisplayName: fmt.Sprintf("User %d", randomN),
		Password:    "abcD1234.",
	}, nil
}
