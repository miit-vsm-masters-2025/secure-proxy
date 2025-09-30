package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var config = createConfig()

func proxyMiddleware(c *gin.Context) {
	host := c.Request.Host
	host = strings.Split(host, ":")[0]
	if host == config.AuthDomain {
		c.Next()
		return
	}
	// TODO proxy
	_ = c.AbortWithError(404, errors.New("Proxying is not implemented yet"))
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.Use(proxyMiddleware)

	r.GET("/set_cookie", func(c *gin.Context) {
		c.SetCookie(
			"SECURE_PROXY_SESSION",
			"E364EEAE-8F50-4B6E-BB9B-E7F56A27160C",
			2592000,
			"/",
			".secure-proxy.wtrn.ru",
			true,
			true,
		)
		c.String(http.StatusOK, "OK")
	})

	r.GET("/proxy", func(c *gin.Context) {
		backendUrl, err := url.Parse("http://127.0.0.1:8000/")
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(backendUrl)
		proxy.Director = func(req *http.Request) {
			req.URL = backendUrl
			req.Header.Set("X-Forwarded-User", "dv.romanov")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	return r
}

func main() {
	router := setupRouter()
	err := router.RunTLS(":8443", "certs/server.pem", "certs/server-key.pem")
	if err != nil {
		panic(err)
	}
}
