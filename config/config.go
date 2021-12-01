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
	GinMode            string
	ServerAddr         string
	ServerPort         string
	StorageAccountName string
	StorageAccountKey  string
	StorageBaseUrl     string
	Shutdown           bool = false
	NoJobWaitTime      int  = 10
)

func InitConfig(file string) error {
	logger.Info("Initalizing configuration")
	loadConfig(file)
	err := configStorage()
	configGin()
	configServer()
	logger.Info("Done initalizing configuration")
	return err
}

func loadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		logger.Error("Could not open env file", err)
		return err
	}
	return nil
}

func configStorage() error {
	var ok bool
	StorageAccountName, ok = os.LookupEnv("STORAGE_ACCOUNT_NAME")
	if !ok || strings.TrimSpace(StorageAccountName) == "" {
		logger.Error("environment variable \"STORAGE_ACCOUNT_NAME\" not set. Cannot start", nil)
		return errors.New("environment variable \"STORAGE_ACCOUNT_NAME\" not set. Cannot start")
	}
	StorageAccountKey, ok = os.LookupEnv("STORAGE_ACCOUNT_KEY")
	if !ok || strings.TrimSpace(StorageAccountKey) == "" {
		logger.Error("environment variable \"STORAGE_ACCOUNT_KEY\" not set. Cannot start", nil)
		return errors.New("environment variable \"STORAGE_ACCOUNT_KEY\" not set. Cannot start")
	}
	StorageBaseUrl, ok = os.LookupEnv("STORAGE_BASE_URL")
	if !ok || strings.TrimSpace(StorageBaseUrl) == "" {
		logger.Error("environment variable \"STORAGE_BASE_URL\" not set. Cannot start", nil)
		return errors.New("environment variable \"STORAGE_BASE_URL\" not set. Cannot start")
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
