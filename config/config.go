package config

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
)

var (
	GinMode    string
	ServerAddr string
	ServerPort string
	DbUser     string
	DbPasswd   string
	DbAddr     string
	DbPort     string
	DbName     string
)

func init() {
	logger.Info("Initalizing configuration")
	loadConfig(".env")
	configGin()
	configServer()
	configDb()
	logger.Info("Done initalizing configuration")
}

func loadConfig(s string) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Could not open env file", err)
	}
}

func configDb() {
	DbUser = os.Getenv("DB_USER")
	DbPasswd = os.Getenv("DB_PASSWD")
	DbAddr = os.Getenv("DB_ADDR")
	DbPort = os.Getenv("DB_PORT")
	DbName = os.Getenv("DB_NAME")
	if strings.TrimSpace(DbUser) == "" ||
		strings.TrimSpace(DbPasswd) == "" ||
		strings.TrimSpace(DbAddr) == "" ||
		strings.TrimSpace(DbPort) == "" ||
		strings.TrimSpace(DbName) == "" {
		logger.Error("DB environment not defined - exiting app", nil)
		panic("DB environment not defined")
	}
}

func configGin() {
	ginMode, ok := os.LookupEnv("GIN_MODE")
	if !ok || (ginMode != gin.ReleaseMode && ginMode != gin.DebugMode && ginMode != gin.TestMode) {
		GinMode = "release"
	}
}

func configServer() {
	var ok bool
	ServerAddr, ok = os.LookupEnv("SERVER_ADDR")
	if !ok {
		ServerAddr = ""
	}
	ServerPort, ok = os.LookupEnv("SERVER_PORT")
	if !ok {
		ServerPort = "8080"
	}
}
