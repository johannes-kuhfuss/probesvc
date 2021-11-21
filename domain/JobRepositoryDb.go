package domain

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type JobRepositoryDb struct {
	client *sqlx.DB
}

func NewJobRepositoryDb(dbClient *sqlx.DB) JobRepositoryDb {
	return JobRepositoryDb{dbClient}
}

func (csdb JobRepositoryDb) FindAll(status string) (*[]Job, api_error.ApiErr) {
	var err error
	jobs := make([]Job, 0)
	if strings.TrimSpace(status) == "" {
		findAllSql := "SELECT * FROM jobList"
		err = csdb.client.Select(&jobs, findAllSql)
	} else {
		findAllSql := "SELECT * FROM jobList WHERE status = ?"
		err = csdb.client.Select(&jobs, findAllSql, status)
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
	err := csdb.client.Get(&job, findByIdSql, id)
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

func (csdb JobRepositoryDb) Create(job Job) api_error.ApiErr {
	insertSql := "INSERT INTO jobs (job_id, name, created_at, created_by, modified_at, modified_by, src_url, status, error_msg, tech_info) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := csdb.client.Exec(insertSql, job.Id, job.Name, job.CreatedAt, job.CreatedBy, job.ModifiedAt, job.ModifiedBy, job.SrcUrl, job.Status, job.ErrorMsg, job.TechInfo)
	if err != nil {
		logger.Error("Error while inserting job into DB", err)
		return api_error.NewInternalServerError("Unexpected database error", nil)
	}
	return nil
}
