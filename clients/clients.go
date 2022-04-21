package clients

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	searchReindex "github.com/ONSdigital/dp-search-reindex-api/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
)

//go:generate moq -out mock/mockZebedeeClient.go -pkg mock . ZebedeeClient
//go:generate moq -out mock/mockSearchReindexClient.go -pkg mock . SearchReindexClient

// ZebedeeClient defines the zebedee client
type ZebedeeClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	GetPublishedIndex(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error)
}

// SearchReindexClient defines the search reindex client
type SearchReindexClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	Health() *health.Client
	PatchJob(context.Context, searchReindexSDK.Headers, string, []searchReindexSDK.PatchOperation) (*searchReindexSDK.RespHeaders, error)
	PostJob(context.Context, searchReindexSDK.Headers) (*searchReindex.Job, error)
	PostTasksCount(ctx context.Context, headers searchReindexSDK.Headers, jobID string, payload []byte) (*searchReindexSDK.RespHeaders, *searchReindex.Task, error)
	GetTask(ctx context.Context, headers searchReindexSDK.Headers, jobID, taskName string) (*searchReindexSDK.RespHeaders, *searchReindex.Task, error)
	URL() string
}
