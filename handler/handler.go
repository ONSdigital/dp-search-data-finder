package handler

import (
	"context"

	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
	dprequest "github.com/ONSdigital/dp-net/request"
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

	//obtain jobid and traceId from the incoming event
	jobid := event.JobID

	//traceid to context
	ctx = dprequest.WithRequestId(ctx, event.TraceID)

	//Make a call to the client passing jobid
	searchReindexTask, err := h.SearchReindexCli.PostTasksCount(ctx, searchReindexSDK.Headers{}, jobid)
	if err != nil {
		return err
	}

    logData = log.Data{
        "Job id": searchReindexTask.JobID,
		"Task Name": searchReindexTask.TaskName,
		"Number of Documents": searchReindexTask.NumberOfDocuments,
    }
	
	log.Info(ctx, "event successfully handled", logData)

	return nil
}
