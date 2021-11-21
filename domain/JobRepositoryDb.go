package domain

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type JobRepositoryDb struct {
	dbclient *sqlx.DB
}

func NewJobRepositoryDb() JobRepositoryDb {
	dbclient, err := sqlx.Open("mysql", "root:Admin000@tcp(192.168.255.128:3306)/jobs?parseTime=true")
	if err != nil {
		panic(err)
	}
	dbclient.SetConnMaxLifetime(time.Minute * 3)
	dbclient.SetMaxOpenConns(10)
	dbclient.SetMaxIdleConns(10)
	return JobRepositoryDb{dbclient: dbclient}
}

func (csdb JobRepositoryDb) FindAll(status string) (*[]Job, api_error.ApiErr) {
	var err error
	jobs := make([]Job, 0)
	if strings.TrimSpace(status) == "" {
		findAllSql := "SELECT * FROM jobList"
		err = csdb.dbclient.Select(&jobs, findAllSql)
	} else {
		findAllSql := "SELECT * FROM jobList WHERE status = ?"
		err = csdb.dbclient.Select(&jobs, findAllSql, status)
	}
	if err != nil {
		logger.Error("Error while querying job list from DB", err)
		return nil, api_error.NewInternalServerError("Unexpected database error", nil)
	}
	if len(jobs) == 0 {
		logger.Info("No data found in DB matching the query")
		return nil, api_error.NewNotFoundError("No data found in DB matching the query")
	}
	return &jobs, nil
}

func (csdb JobRepositoryDb) FindById(id string) (*Job, api_error.ApiErr) {
	var job Job
	findByIdSql := "SELECT * FROM jobList WHERE job_id = ?"
	err := csdb.dbclient.Get(&job, findByIdSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api_error.NewNotFoundError(fmt.Sprintf("Job with id %v not found", id))
		} else {
			logger.Error("Error while scanning job from DB", err)
			return nil, api_error.NewInternalServerError("Unexpected database error", nil)
		}
	}
	return &job, nil
}
