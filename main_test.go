package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
)

func Test_generateKey(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "alice@example.com",
	})
	if err != nil {
		panic(err)
	}

	secret := key.Secret()
	println("Secret: ", secret)
}

func Test_validateTotp(t *testing.T) {
	secret := "DU2DGNS3ALWXIGLWK7ZVBXHL7ZN6ZCAC"
	generatedCode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		panic(err)
	}
	println("Code: ", generatedCode)

	validationResult := totp.Validate(generatedCode, secret)
	println("Validation passed: ", validationResult)
}

func Test_valkey(t *testing.T) {
	host := "localhost"
	port := 6379

	config := config.NewClientConfiguration().
		WithAddress(&config.NodeAddress{Host: host, Port: port})

	client, err := glide.NewClient(config)
	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}

	res, err := client.Ping(context.Background())
	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}
	fmt.Println(res) // PONG

	client.Close()

}
