package clients

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-search-reindex-api/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
)

//go:generate moq -out mock/reindexapi.go -pkg mock . SearchReindexClient

// DatasetApiClient defines the zebedee client
type SearchReindexClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	PostTasksCount(ctx context.Context, headers searchReindexSDK.Headers, jobID string, payload []byte) (models.Task, error)
}
