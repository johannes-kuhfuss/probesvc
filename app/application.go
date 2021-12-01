package app

import (
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/handler"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	router      *gin.Engine
	jobHandler  handler.JobHandlers
	azureClient *azblob.ServiceClient
	jobService  service.DefaultJobService
	fileService service.DefaultFileService
)

func connectToAzureBlob() (*azblob.ServiceClient, api_error.ApiErr) {
	blobUrl, err := url.Parse(config.StorageBaseUrl)
	if err != nil {
		logger.Error("Cannot parse storage base URL", nil)
		return nil, api_error.NewBadRequestError("Cannot parse storage base URL")
	}

	cred, err := azblob.NewSharedKeyCredential(config.StorageAccountName, config.StorageAccountKey)
	if err != nil {
		logger.Error("Cannot access storage account - wrong credentials", err)
		return nil, api_error.NewInternalServerError("Cannot access storage account - wrong credentials", err)
	}

	serviceClient, err := azblob.NewServiceClientWithSharedKey(blobUrl.String(), cred, nil)

	if err != nil {
		logger.Error("Cannot access storage account - could not create service client", err)
		return nil, api_error.NewInternalServerError("Cannot access storage account - could not create service client", err)
	}
	return &serviceClient, nil
}

func initRouter() {
	gin.SetMode(config.GinMode)
	gin.DefaultWriter = logger.GetLogger()
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}

func wireApp() {
	customerRepo := domain.NewJobRepositoryMem()
	jobService = service.NewJobService(customerRepo)
	jobHandler = handler.JobHandlers{Service: jobService}
	azureFileRepo := domain.NewFileRepositoryAzure(azureClient)
	fileService = service.NewFileService(azureFileRepo)
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
		panic(err)
	}
	azureClient, err = connectToAzureBlob()
	if err != nil {
		panic(err)
	}
	initRouter()
	wireApp()
	mapUrls()
	startRouter()
	startProcessing()
	logger.Info("Application ended")
}

func startProcessing() {
	go fileService.Run()
}
