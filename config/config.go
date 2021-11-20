package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	ginMode    = "debug" // release, debug, test
	ListenAddr = ""
)

func init() {
	logger.Info("Initalizing configuration")
	osGinMode := os.Getenv("GIN_MODE")
	if osGinMode == gin.ReleaseMode || osGinMode == gin.DebugMode || osGinMode == gin.TestMode {
		ginMode = osGinMode
	}
	logger.Debug(fmt.Sprintf("Gin-Gonic Mode: %v\n", ginMode))
	ListenAddr = os.Getenv("LISTEN_ADDR")
	if len(ListenAddr) == 0 {
		ListenAddr = ":8080"
	}
	logger.Info("Done initalizing configuration")
}

func GinMode() string {
	return ginMode
}
