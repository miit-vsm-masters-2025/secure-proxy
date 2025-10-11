package main

import (
	"secure-proxy/impl"
)

func main() {
	router := impl.SetupRouter()
	err := router.RunTLS(":8443", "certs/server.pem", "certs/server-key.pem")
	if err != nil {
		panic(err)
	}
}
