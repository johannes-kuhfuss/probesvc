package domain

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

//go:generate mockgen -destination=../mocks/domain/mockFileRepository.go -package=domain github.com/johannes-kuhfuss/probesvc/domain FileRepository
type FileRepository interface {
	GetClient() *azblob.ServiceClient
}
