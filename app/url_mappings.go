package app

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/handler"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
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

func mapUrls() {
	logger.Debug("Mapping URLs")
	dbclient := getDbClient()
	jh := handler.JobHandlers{Service: service.NewJobService(domain.NewJobRepositoryDb(dbclient))}
	router.GET("/jobs", jh.GetAllJobs)
	router.GET("jobs/:job_id", jh.GetJobById)
	logger.Debug("Done mapping URLs")
}
