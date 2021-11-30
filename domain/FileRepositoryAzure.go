package domain

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type FileRepositoryAzure struct {
	serviceClient *azblob.ServiceClient
}

func NewFileRepositoryAzure() FileRepositoryAzure {
	client, err := connectToAzureBlob()
	if err != nil {
		logger.Error("Ouch.", err)
	}
	return FileRepositoryAzure{client}
}

func connectToAzureBlob() (*azblob.ServiceClient, api_error.ApiErr) {
	blobUrl, err := url.Parse(config.StorageBaseUrl)
	if err != nil {
		logger.Error("Cannot parse source URL", nil)
		return nil, api_error.NewBadRequestError("Cannot parse source URL")
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
