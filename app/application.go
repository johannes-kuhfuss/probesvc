package app

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/handler"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	router     *gin.Engine
	dbClient   *sqlx.DB
	jobHandler handler.JobHandlers
)

func getDbClient() *sqlx.DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DbUser, config.DbPasswd, config.DbAddr, config.DbPort, config.DbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}

func initRouter() {
	gin.SetMode(config.GinMode)
	gin.DefaultWriter = logger.GetLogger()
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}

func wireApp() {
	dbClient = getDbClient()
	customerRepositoryDb := domain.NewJobRepositoryDb(dbClient)
	jobHandler = handler.JobHandlers{Service: service.NewJobService(customerRepositoryDb)}
}

func startRouter() {
	listenAddr := fmt.Sprintf("%s:%s", config.ServerAddr, config.ServerPort)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := router.Run(listenAddr); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile)
	if err != nil {
		panic("Error while configuring the application")
	}
	initRouter()
	wireApp()
	mapUrls()
	startRouter()
	logger.Info("Application ended")
}
