package handler

import (
	"context"

	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/log.go/v2/log"
)

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeCli clients.ZebedeeClient
	Config     config.Config
}

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, event *models.ReindexRequested) error {
	logData := log.Data{
		"event": event,
	}

	log.Info(ctx, "reindex requested event handler called", logData)

	publishedIndex, err := h.ZebedeeCli.GetPublishedIndex(ctx, nil)
	if err != nil {
		// the cfg.ZebedeeClientTimeout maybe needs to be increased
		return err
	}

	urlList := make([]string, 10)
	publishedItems := publishedIndex.Items
	totalZebedeeDocs := len(publishedItems)

	// only the first 10 docs are retrieved for testing and performance purposes
	// it takes more than 10 mins to retrieve all document urls from zebedee
	// TODO: remove (i < 10) condition when this app has been completely implemented
	for i := 0; (i < 10) && (i < totalZebedeeDocs); i++ {
		urlList[i] = publishedItems[i].URI
	}

	log.Info(ctx, "first 10 Zebedee docs URLs retrieved", log.Data{"first URLs": urlList})
	log.Info(ctx, "total number of zebedee docs retrieved", log.Data{"total_documents": totalZebedeeDocs})

	log.Info(ctx, "event successfully handled", logData)
	return nil
}
