package main

import (
	"context"
	"fmt"
	"log"
	"merch_service/internal/client"
)

func main() {
	// Демо версия. Позже уберем
	user := client.Credentials{
		Login: "aboba",
		Pass:  "123123",
	}

	cli := client.NewClient()

	err := cli.Register(context.Background(), &user)

	if err != nil {
		log.Fatalln(err)
	}

	tokens, err := cli.GetTokens(context.Background(), &user)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Ваш токен", tokens.Token)
}
