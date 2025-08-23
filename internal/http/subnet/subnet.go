package subnet

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewTrustedSubnetMiddleware(trustedSubnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if trustedSubnet == "" {
			c.Next()
			return
		}

		realIP := c.GetHeader("X-Real-IP")
		if realIP == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		_, subnet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ip := net.ParseIP(realIP)
		if ip == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !subnet.Contains(ip) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
