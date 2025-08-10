package api

import (
	"github.com/gin-gonic/gin"
	"net"
	"picp/config"
	"picp/logger"
)

var engine *gin.Engine

func Run(server net.Listener) error {
	engine = gin.New()
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logger.GetWitter(),
	}), gin.Recovery())
	if config.Common.User != "" && config.Common.Password != "" {
		engine.Use(gin.BasicAuth(gin.Accounts{
			config.Common.User: config.Common.Password,
		}))
	}
	initPages()
	initApi(engine.Group("/api"))
	return engine.RunListener(server)
}
