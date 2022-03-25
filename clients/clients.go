package clients

import (
	"context"

	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	searchReindex "github.com/ONSdigital/dp-search-reindex-api/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
)

//go:generate moq -out mock/clients.go -pkg mock . ZebedeeClient

// ZebedeeClient defines the zebedee client
type ZebedeeClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	GetPublishedIndex(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error)
}

// SearchReindexClient defines the search reindex client
type SearchReindexClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	PostJob(context.Context, searchReindexSDK.Headers) (searchReindex.Job, error)
}
