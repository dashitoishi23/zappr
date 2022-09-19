package azurestorageservice

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type AzureStorageService interface {
	CreateAzureStorageContainer(ctx context.Context, containerName string, isPubliclyAccessible bool) (bool, error)
	DeleteAzureStorageContainer(ctx context.Context, containerName string) (bool, error)
}

type azureStorageService struct {
	ServiceURL azblob.ServiceURL
}

func NewAzureStorageService() (AzureStorageService, error) {
	accountName := os.Getenv("AZURE_STORAGE_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return &azureStorageService{}, err
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	storageAccountURL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	if err != nil {
		return &azureStorageService{}, err
	}

	return &azureStorageService{
		ServiceURL: azblob.NewServiceURL(*storageAccountURL, pipeline),
	}, nil
}

func (a *azureStorageService) CreateAzureStorageContainer(ctx context.Context, containerName string, 
	isPubliclyAccessible bool) (bool, error) {
	containerURL := a.ServiceURL.NewContainerURL(containerName)

	var publicAccessType azblob.PublicAccessType

	if isPubliclyAccessible {
		publicAccessType = azblob.PublicAccessContainer
	} else {
		publicAccessType = azblob.PublicAccessNone
	}


	_, err := containerURL.Create(ctx, azblob.Metadata{}, publicAccessType)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *azureStorageService) DeleteAzureStorageContainer(ctx context.Context, containerName string) (bool, error) {	
	containerURL := a.ServiceURL.NewContainerURL(containerName)

	_, err := containerURL.Delete(ctx, azblob.ContainerAccessConditions{})

	if err != nil {
		return false, err

	}
	return true, nil
}
