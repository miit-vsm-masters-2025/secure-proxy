package impl

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.Use(errorHandlingMiddleware)
	r.Use(proxyMiddleware)
	r.LoadHTMLGlob("templates/*")

	r.GET("/", renderAuthPage)
	r.POST("/auth", validateTotp)

	// Методы ниже нам на самом деле не нужны, но пока оставляю их для примера. Нужно будет удалить после завершения реализации.
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

func renderAuthPage(c *gin.Context) {
	redirectUrl := c.Query("redirectUrl")
	c.HTML(http.StatusOK, "auth.tmpl", gin.H{
		"redirectUrl": redirectUrl,
	})
	// Отрендерить и вернуть html-страницу с формой для ввода логина и TOTP-кода
}

func validateTotp(c *gin.Context) {
	username := c.PostForm("username")
	enteredTotpCode := c.PostForm("totp")
	redirectUrl := c.PostForm("redirectUrl")

	user, err := findUserByUsername(username)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	secret := user.TotpSecret
	// TODO: check totp code
	totpOk := totp.Validate(enteredTotpCode, secret)
	if !totpOk {
		c.AbortWithError(404, fmt.Errorf("incorrect totp code for user %s", user.Username))
		return
	}

	uuid := uuid.New().String()
	err = valkeyClient.createSession(c, username, uuid)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to create new session in valkey: %s", err))
		return
	}

	c.SetCookie(
		config.CookieName,
		uuid,
		2592000, // TODO: We probably should refresh cookie in middleware
		"/",
		config.Sessions.CookieDomain,
		true,
		true,
	)

	c.Redirect(302, redirectUrl)

	// Проверить введенный TOTP. Если все ок - проставить куку и отредиректить на url, указанный в параметре redirectUrl.
	// Если нет - отрендерить ту же форму что и в методе выше, но с сообщением об ошибке.
}

func findUserByUsername(username string) (User, error) {
	for i := range config.Users {
		user := config.Users[i]
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user %s not found", username)
}
