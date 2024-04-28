package main

import (
	"fmt"
	"time"
	"tokenvault"
)

func main() {

	duration := time.Second * 5

	tokenManager := tokenvault.NewTokenManager(
		"RANDOM_API_TOKEN",
		&duration,
		func() (interface{}, error) {
			// some api call should be replaced by the generator function
			// for now just using a time.Sleep to simulate an api call
			time.Sleep(time.Second * 1)

			token := "token" + fmt.Sprint(time.Now().Unix())

			return token, nil
		},
	)

	go tokenManager.RunTokenGenerator() // running the generator in the background which will update the token

	for {
		token := tokenManager.GetToken()
		fmt.Println("token", token)
		time.Sleep(time.Second * 6)
	}
}
