package event

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// ReindexRequestedHandler is the hander for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeClient *zebedee.Client
}

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, cfg *config.Config, event *ReindexRequested) (err error) {
	logData := log.Data{
		"event": event,
	}
	log.Info(ctx, "reindex requested event handler called", logData)

	publishedIndex, err := h.ZebedeeClient.GetPublishedIndex(ctx, nil)

	publishedItem := publishedIndex.Items[0]

	//TODO Log out the first 10 URLs retrieved
	log.Info(ctx, "URL retrieved: " + publishedItem.URI)
	log.Info(ctx, "event successfully handled", logData)

	return nil
}
