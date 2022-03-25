package handler

import (
	"context"

	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/models"
	searchReindex "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeCli       clients.ZebedeeClient
	SearchReindexCli clients.SearchReindexClient
}

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, event *models.ReindexRequested) (err error) {
	logData := log.Data{
		"event": event,
	}
	log.Info(ctx, "reindex requested event handler called", logData)

	if h.ZebedeeCli == nil {
		return errors.New("the Zebedee client in the reindex requested handler must not be nil")
	}

	publishedIndex, err := h.ZebedeeCli.GetPublishedIndex(ctx, nil)
	if err != nil {
		return err
	}

	urlList := make([]string, 10)
	publishedItems := publishedIndex.Items
	for i := 0; (i < 10) && (i < len(publishedItems)); i++ {
		urlList[i] = publishedItems[i].URI
	}
	log.Info(ctx, "first 10 URLs retrieved", log.Data{"first URLs": urlList})
	log.Info(ctx, "event successfully handled", logData)

    if h.SearchReindexCli == nil {
        return errors.New("the search reindex client in the reindex requested handler must not be nil")
    }

    searchReindexJob, err := h.SearchReindexCli.PostJob(ctx, searchReindex.Headers{})
    if err != nil {
        return err
    }

    logJob := log.Data{
        "new job": searchReindexJob,
    }

    log.Info(ctx, "new job created", logJob)

	return nil
}
