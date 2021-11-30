package domain

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type FileRepositoryAzure struct {
	serviceClient *azblob.ServiceClient
}

func NewFileRepositoryAzure(client *azblob.ServiceClient) FileRepositoryAzure {
	return FileRepositoryAzure{client}
}
