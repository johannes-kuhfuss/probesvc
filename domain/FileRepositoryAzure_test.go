package domain

import (
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/stretchr/testify/assert"
)

var (
	fileRepo    FileRepositoryAzure
	azureClient *azblob.ServiceClient
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

func setupFile() func() {
	azureClient, err := connectToAzureBlob()
	if err != nil {
		panic(err)
	}
	fileRepo = NewFileRepositoryAzure(azureClient)
	return func() {
		fileRepo.serviceClient = nil
		azureClient = nil
	}
}

func Test_GetClient_Returns_AzureClient(t *testing.T) {
	teardown := setupFile()
	defer teardown()
	myClient := fileRepo.GetClient()

	assert.NotNil(t, myClient)
	assert.IsType(t, azureClient, myClient)
}
