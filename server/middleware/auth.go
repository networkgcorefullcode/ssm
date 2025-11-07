package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/networkgcorefullcode/ssm/const"
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

		// check if the user is valid
		if jwtPayload.Sub != constants.USER_UDM && jwtPayload.Sub != constants.USER_WEBCONSOLE {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}

		// check if the operation is allow for the role (udm only decrypt, webconsole all actions in the list)
		action := determineAction(c)
		if !slices.Contains(constants.ActionList, action) && jwtPayload.Sub == constants.USER_WEBCONSOLE {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid operation for the user"})
			return
		}

		if action != constants.ACTION_DECRYPT_DATA && action != constants.ACTION_HEALTH_CHECK && jwtPayload.Sub == constants.USER_UDM {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid operation for the user"})
			return
		}

		// Set the payload in context for further handlers
		c.Set("jwt-payload", jwtPayload)

		// If all is well, proceed to the handler
		c.Next()
	}
}
