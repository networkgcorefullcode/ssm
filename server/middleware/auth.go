package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

func AuthenticateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.GetHeader("Authorization")

		if jwtToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth headers"})
			return
		}

		tokenString := strings.Replace(jwtToken, "Bearer ", "", 1)

		session := mgr.GetSession()
		defer mgr.LogoutSession(session)

		// verify JWT token here
		jwtPayload, err := pkcs11mgr.VerifyJWT(session, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Set the payload in context for further handlers
		c.Set("jwt-payload", jwtPayload)

		// If all is well, proceed to the handler
		c.Next()
	}
}
