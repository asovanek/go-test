package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS configures cross-origin responses for browser clients.
func CORS(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.Writer.Header()
		hdr.Set("Access-Control-Allow-Origin", origin)
		hdr.Set("Access-Control-Allow-Methods", strings.Join([]string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
		}, ", "))
		hdr.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		hdr.Set("Access-Control-Expose-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
