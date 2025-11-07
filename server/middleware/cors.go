package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
)

// configureCORS sets up CORS middleware based on configuration settings yaml file
func ConfigureCORS(r *gin.Engine) {
	corsConfig := factory.SsmConfig.Configuration.CORS

	logger.AppLog.Infof("Configuring CORS middleware, CORS: %+v", corsConfig)

	config := cors.Config{
		AllowAllOrigins:           getBoolValue(corsConfig.AllowAllOrigins),
		AllowOrigins:              corsConfig.AllowOrigins,
		AllowMethods:              corsConfig.AllowMethods,
		AllowHeaders:              corsConfig.AllowHeaders,
		AllowCredentials:          getBoolValue(corsConfig.AllowCredentials),
		ExposeHeaders:             corsConfig.ExposeHeaders,
		MaxAge:                    corsConfig.GetMaxAgeDuration(),
		AllowWildcard:             getBoolValue(corsConfig.AllowWildcard),
		AllowBrowserExtensions:    getBoolValue(corsConfig.AllowBrowserExtensions),
		AllowWebSockets:           getBoolValue(corsConfig.AllowWebSockets),
		AllowFiles:                getBoolValue(corsConfig.AllowFiles),
		AllowPrivateNetwork:       getBoolValue(corsConfig.AllowPrivateNetwork),
		CustomSchemas:             corsConfig.CustomSchemas,
		OptionsResponseStatusCode: corsConfig.GetOptionsStatusCode(),
	}

	r.Use(cors.New(config))
}

// getBoolValue safely dereferences a bool pointer, returning false if nil
func getBoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
