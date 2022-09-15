package azurestorage

import (
	"context"
)

type AzureStorageService interface {
	CreateAzureStorageContainer(ctx context.Context, containerName string) (bool, error)
	DeleteAzureStorageContainer(ctx context.Context, containerName string) (bool, error)
}

type azureStorageService struct {

}

func NewAzureStorageService() AzureStorageService {
	return &azureStorageService{} 
}