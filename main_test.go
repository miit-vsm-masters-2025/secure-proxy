package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/valkey-io/valkey-go"
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
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		panic(err)
	}

	context := context.Background()
	resp, err := client.Do(context, client.B().Ping().Build()).ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp) // PONG

	key := "session_E364EEAE-8F50-4B6E-BB9B-E7F56A27160C"
	value := "dv.romanov"

	err = client.Do(context, client.B().Setex().Key(key).Seconds(5).Value(value).Build()).Error()
	if err != nil {
		panic(err)
	}

	retrieved, err := client.Do(context, client.B().Getex().Key(key).ExSeconds(60).Build()).ToString()
	if err != nil {
		panic(err)
	}

	println("retrieved: ", retrieved)

	client.Close()

}
