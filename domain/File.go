package domain

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

type FileRepository interface {
	GetClient() *azblob.ServiceClient
}
