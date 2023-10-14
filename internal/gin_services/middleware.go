package ginservices

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware function to restrict access to IP addresses within the local network
func restrictToLocalNetwork() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client IP address
		clientIP := c.ClientIP()

		// Parse the client IP address
		ip := net.ParseIP(clientIP)

		// Define the local network IP ranges as strings
		localNetworks := []string{
			"192.168.166.0/24",
			"127.0.0.1/32",
		}

		// Check if the client IP address is within any of the local network ranges
		for _, network := range localNetworks {
			_, ipNet, err := net.ParseCIDR(network)
			if err != nil {
				// Invalid network range, skip to the next one
				continue
			}

			if ipNet.Contains(ip) {
				// Allow access to the API if the client IP address is within any of the local network ranges
				c.Next()
				return
			}
		}

		// Deny access to the API if the client IP address is not within any of the local network ranges
		c.AbortWithStatus(http.StatusForbidden)
	}
}
