package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

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

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
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
