package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSOptions configures the CORS middleware.
type CORSOptions struct {
	AllowedOrigins   []string
	AllowLocalhost   bool // allow http://localhost:* and http://127.0.0.1:* (dev)
}

// CORS configures cross-origin responses for browser clients.
func CORS(opts CORSOptions) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(opts.AllowedOrigins))
	for _, origin := range opts.AllowedOrigins {
		allowed[origin] = struct{}{}
	}

	allowOrigin := func(origin string) bool {
		if origin == "" {
			return false
		}
		if _, ok := allowed[origin]; ok {
			return true
		}
		return opts.AllowLocalhost && isLocalhostOrigin(origin)
	}

	return func(c *gin.Context) {
		hdr := c.Writer.Header()
		origin := c.GetHeader("Origin")
		if allowOrigin(origin) {
			hdr.Set("Access-Control-Allow-Origin", origin)
			hdr.Set("Vary", "Origin")
		}

		hdr.Set("Access-Control-Allow-Methods", strings.Join([]string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
		}, ", "))
		hdr.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		hdr.Set("Access-Control-Expose-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			if origin != "" && !allowOrigin(origin) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func isLocalhostOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" {
		return false
	}
	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1"
}
