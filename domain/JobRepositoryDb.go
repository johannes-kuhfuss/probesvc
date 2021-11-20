package domain

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"

	_ "github.com/go-sql-driver/mysql"
)

type JobRepositoryDb struct {
	dbclient *sql.DB
}

func NewJobRepositoryDb() JobRepositoryDb {
	dbclient, err := sql.Open("mysql", "root:Admin000@tcp(192.168.255.128:3306)/jobs")
	if err != nil {
		panic(err)
	}
	dbclient.SetConnMaxLifetime(time.Minute * 3)
	dbclient.SetMaxOpenConns(10)
	dbclient.SetMaxIdleConns(10)
	return JobRepositoryDb{dbclient: dbclient}
}

func (csdb JobRepositoryDb) FindAll() (*Jobs, api_error.ApiErr) {
	findAllSql := "SELECT job_id, name, src_url, status FROM jobList"
	rows, err := csdb.dbclient.Query(findAllSql)
	if err != nil {
		logger.Error("Error while querying job list from DB", err)
		return nil, api_error.NewInternalServerError("Unexpected database error", nil)
	}
	jobs := make(Jobs, 0)
	for rows.Next() {
		var j Job
		err := rows.Scan(&j.Id, &j.Name, &j.SrcUrl, &j.Status)
		if err != nil {
			logger.Error("Error while scanning job list from DB", err)
			return nil, api_error.NewInternalServerError("Unexpected database error", nil)
		}
		jobs = append(jobs, j)
	}
	return &jobs, nil
}

func (csdb JobRepositoryDb) FindById(id string) (*Job, api_error.ApiErr) {
	findByIdSql := "SELECT job_id, name, src_url, status FROM jobList WHERE job_id = ?"
	row := csdb.dbclient.QueryRow(findByIdSql, id)
	var j Job
	err := row.Scan(&j.Id, &j.Name, &j.SrcUrl, &j.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api_error.NewNotFoundError(fmt.Sprintf("customer with id %v not found", id))
		} else {
			logger.Error("Error while scanning job from DB", err)
			return nil, api_error.NewInternalServerError("Unexpected database error", nil)
		}
	}
	return &j, nil
}
