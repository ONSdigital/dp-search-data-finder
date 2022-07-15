package handler

import (
	"context"

	dprequest "github.com/ONSdigital/dp-net/request"
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

type key int

const keyRequestID key = iota

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, event *models.ReindexRequested) error {
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

	// add trace id to context
	ctx = dprequest.WithRequestId(ctx, event.TraceID)

	log.Info(ctx, "event successfully handled", logData)
	return nil
}
