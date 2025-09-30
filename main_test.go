package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
	"github.com/valkey-io/valkey-glide/go/v2/constants"
	"github.com/valkey-io/valkey-glide/go/v2/options"
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

	context := context.Background()
	res, err := client.Ping(context)
	if err != nil {
		fmt.Println("There was an error: ", err)
		return
	}
	fmt.Println(res) // PONG

	key := "session_E364EEAE-8F50-4B6E-BB9B-E7F56A27160C"
	value := "dv.romanov"

	_, err = client.SetWithOptions(context, key, value, options.SetOptions{
		Expiry: &options.Expiry{
			Type:     constants.Seconds,
			Duration: 5,
		},
	})
	if err != nil {
		panic(err)
	}

	retrieved, err := client.Get(context, key)
	if err != nil {
		panic(err)
	}

	println("retrieved: ", retrieved.Value())

	client.Close()

}
