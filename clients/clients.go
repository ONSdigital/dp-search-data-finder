package clients

import (
	"context"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mock/mockZebedeeClient.go -pkg mock . ZebedeeClient
//go:generate moq -out mock/mockDatasetAPIClient.go -pkg mock . DatasetAPIClient

// ZebedeeClient defines the zebedee client
type ZebedeeClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	GetPublishedIndex(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error)
}

// DatasetAPIClient defines the dataset-api client
type DatasetAPIClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	GetDatasets(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string, q *datasetclient.QueryParams) (m datasetclient.List, err error)
	GetFullEditionsDetails(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m []datasetclient.EditionsDetails, err error)
	GetVersionMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m datasetclient.Metadata, err error)
}
