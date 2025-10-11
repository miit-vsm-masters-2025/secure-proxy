package impl

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func proxyRequest(upstream *Upstream, c *gin.Context) {
	upstreamUrl, err := url.Parse(upstream.Destination + c.Request.RequestURI)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(upstreamUrl)
	proxy.Director = func(req *http.Request) {
		req.URL = upstreamUrl
		req.Header.Set("X-Forwarded-User", "dv.romanov")
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
