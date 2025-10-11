package impl

import (
	"net/http"
	"net/url"

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

	username := getUsername(c)
	if username == "" {
		redirectToAuth(c)
		return
	}
	// TODO proxy
	for i := range config.Upstreams {
		upstream := &config.Upstreams[i]
		if upstream.Host == host {
			proxyRequest(upstream, c)
			return
		}
	}

	c.String(http.StatusNotFound, "Upstream Not Found")
}

func getUsername(c *gin.Context) string {
	// TODO: Get cookie, find username in valkey, if not found - redirect to auth
	// Use GetExWithOptions to extend key expiration time
	return ""
}

func redirectToAuth(c *gin.Context) {
	authenticatedRedirectUrl := "https://" + c.Request.Host + c.Request.RequestURI
	authUrl := "https://" + config.AuthDomain + "/?redirectUrl=" + url.QueryEscape(authenticatedRedirectUrl)
	c.Redirect(http.StatusFound, authUrl)
}
