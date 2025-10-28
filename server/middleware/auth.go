package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var secretStore map[string]string = map[string]string{
	"udm":        "",
	"webconsole": "",
}

func AuthenticateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID := c.GetHeader("X-Service-Id")
		timestamp := c.GetHeader("X-Timestamp")
		signature := c.GetHeader("X-Signature")

		if serviceID == "" || timestamp == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth headers"})
			return
		}

		secret, err := getServiceSecret(serviceID) // get the secret from database MongoDB
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown service"})
			return
		}

		// Verify that the timestamp is not out of range (prevents replays)
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || time.Since(time.Unix(ts, 0)) > 2*time.Minute {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired timestamp"})
			return
		}

		// Recalculate the expected signature
		data := c.Request.Method + ":" + c.Request.URL.Path + ":" + timestamp
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(data))
		expectedSig := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(expectedSig), []byte(signature)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		// If all is well, proceed to the handler
		c.Next()
	}
}
