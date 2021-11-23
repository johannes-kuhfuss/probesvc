package config

import (
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
)

const (
	EnvFile = ".env"
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

func InitConfig(file string) error {
	logger.Info("Initalizing configuration")
	loadConfig(file)
	configGin()
	configServer()
	err := configDb()
	if err != nil {
		return err
	}
	logger.Info("Done initalizing configuration")
	return nil
}

func loadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		logger.Error("Could not open env file", err)
		return err
	}
	return nil
}

func configDb() error {
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
		return errors.New("DB environment not defined")
	}
	return nil
}

func configGin() {
	ginMode, ok := os.LookupEnv("GIN_MODE")
	if !ok || (ginMode != gin.ReleaseMode && ginMode != gin.DebugMode && ginMode != gin.TestMode) {
		GinMode = "release"
	} else {
		GinMode = ginMode
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
