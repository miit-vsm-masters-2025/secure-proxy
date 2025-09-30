package main

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
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
