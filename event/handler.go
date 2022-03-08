package event

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/log.go/v2/log"
)

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeClient *zebedee.Client
}

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, event *ReindexRequested) (err error) {
	logData := log.Data{
		"event": event,
	}
	log.Info(ctx, "reindex requested event handler called", logData)

	publishedIndex, err := h.ZebedeeClient.GetPublishedIndex(ctx, nil)
	if err != nil {
		return err
	}

	urlList := make([]string, 10)
	publishedItems := publishedIndex.Items
	for i := 0; i < 3; i++ {
		urlList[i] = publishedItems[i].URI
	}
	log.Info(ctx, "First 10 URLs retrieved", log.Data{"first URLs": urlList})
	log.Info(ctx, "event successfully handled", logData)

	return nil
}
