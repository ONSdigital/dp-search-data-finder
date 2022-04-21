package handler

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
)

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeCli       clients.ZebedeeClient
	SearchReindexCli clients.SearchReindexClient
	Config           config.Config
}

type key int

const keyRequestID key = iota

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, event *models.ReindexRequested) (err error) {
	logData := log.Data{
		"event": event,
	}
	ctx = context.WithValue(ctx, keyRequestID, event.TraceID)
	log.Info(ctx, "reindex requested event handler called", logData)

	publishedIndex, err := h.ZebedeeCli.GetPublishedIndex(ctx, nil)
	if err != nil {
		return err
	}

	urlList := make([]string, 10)
	publishedItems := publishedIndex.Items
	totalZebedeeDocs := len(publishedItems)
	for i := 0; (i < 10) && (i < totalZebedeeDocs); i++ {
		urlList[i] = publishedItems[i].URI
	}
	log.Info(ctx, "first 10 URLs retrieved", log.Data{"first URLs": urlList})

	if h.SearchReindexCli == nil {
		return errors.New("the search reindex client in the reindex requested handler must not be nil")
	}

	headers := searchReindexSDK.Headers{
		IfMatch:          "*",
		ServiceAuthToken: h.Config.ServiceAuthToken,
	}

	patchList := make([]searchReindexSDK.PatchOperation, 2)
	statusOperation := searchReindexSDK.PatchOperation{
		Op:    "replace",
		Path:  "/state",
		Value: "in-progress",
	}
	patchList[0] = statusOperation
	totalDocsOperation := searchReindexSDK.PatchOperation{
		Op:    "replace",
		Path:  "/total_search_documents",
		Value: totalZebedeeDocs,
	}
	patchList[1] = totalDocsOperation

	logPatchList := log.Data{
		"patch list": patchList,
	}

	log.Info(ctx, "patch list for request", logPatchList)

	var respHeaders *searchReindexSDK.RespHeaders
	respHeaders, err = h.SearchReindexCli.PatchJob(ctx, headers, event.JobID, patchList)
	if err != nil {
		return err
	}
	respETag := respHeaders.ETag

	log.Info(ctx, "eTag returned in response:"+respETag)

	log.Info(ctx, "event successfully handled", logData)
	return nil
}
