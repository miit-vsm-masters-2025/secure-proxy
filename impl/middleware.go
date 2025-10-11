package impl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func proxyMiddleware(c *gin.Context) {
	host := c.Request.Host
	//host = strings.Split(host, ":")[0]
	if host == config.AuthDomain {
		c.Next()
		return
	}

	c.Abort()

	username, err := getUsername(c)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	if username == "" {
		redirectToAuth(c)
		return
	}
	// TODO proxy
	host = strings.Split(host, ":")[0] // Remove port from host
	for i := range config.Upstreams {
		upstream := &config.Upstreams[i]
		if upstream.Host == host {
			proxyRequest(upstream, c)
			return
		}
	}

	c.String(http.StatusNotFound, "Upstream Not Found")
}

func getUsername(c *gin.Context) (string, error) {
	sessionKey, err := c.Cookie(config.CookieName)
	if err != nil {
		return "", nil
	}

	username, err := valkeyClient.findUsernameBySession(c, sessionKey)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve username from vault: %s", err)
	}

	return username, nil
}

func redirectToAuth(c *gin.Context) {
	authenticatedRedirectUrl := "https://" + c.Request.Host + c.Request.RequestURI
	authUrl := "https://" + config.AuthDomain + "/?redirectUrl=" + url.QueryEscape(authenticatedRedirectUrl)
	c.Redirect(http.StatusFound, authUrl)
}
