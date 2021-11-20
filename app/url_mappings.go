package app

import "github.com/johannes-kuhfuss/services_utils/logger"

func mapUrls() {
	logger.Debug("Mapping URLs")

	/*
		router.GET("/ping", controllers.PingController.Pong)
		router.POST("/job", controllers.JobController.Create)
		router.GET("/job/:job_id", controllers.JobController.Get)
		router.DELETE("/job/:job_id", controllers.JobController.Delete)
		router.PUT("/job/:job_id", controllers.JobController.Update)
		router.PATCH("/job/:job_id", controllers.JobController.UpdatePart)
		router.GET("/jobs/", controllers.JobController.GetAll)
	*/

	logger.Debug("Done mapping URLs")
}
