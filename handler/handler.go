package handler

import (
	"context"
	"errors"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/models"
	searchReindex "github.com/ONSdigital/dp-search-reindex-api/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	zebedeeTaskName = "zebedee"
)

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeCli       clients.ZebedeeClient
	SearchReindexCli searchReindexSDK.Client
	Config           config.Config
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

	// obtain jobid
	jobID := event.JobID

	// traceid to context
	ctx = dprequest.WithRequestId(ctx, event.TraceID)

	if h.SearchReindexCli == nil {
		return errors.New("the search reindex client in the reindex requested handler must not be nil")
	}

	headers := searchReindexSDK.Headers{
		IfMatch:          "*",
		ServiceAuthToken: h.Config.ServiceAuthToken,
	}

	patchList := make([]searchReindexSDK.PatchOperation, 2)
	patchList[0] = searchReindexSDK.PatchOperation{
		Op:    "replace",
		Path:  "/state",
		Value: "in-progress",
	}
	patchList[1] = searchReindexSDK.PatchOperation{
		Op:    "replace",
		Path:  "/total_search_documents",
		Value: totalZebedeeDocs,
	}

	logPatchList := log.Data{
		"patch list": patchList,
	}

	log.Info(ctx, "patch list for request", logPatchList)

	var respHeaders *searchReindexSDK.RespHeaders
	respHeaders, err = h.SearchReindexCli.PatchJob(ctx, headers, jobID, patchList)
	if err != nil {
		return err
	}
	patchJobRespETag := respHeaders.ETag

	log.Info(ctx, "job state and total_search_documents were patched successfully")

	taskToCreate := searchReindex.TaskToCreate{
		TaskName:          zebedeeTaskName,
		NumberOfDocuments: totalZebedeeDocs,
	}

	var searchReindexTask *searchReindex.Task
	respHeaders, searchReindexTask, err = h.SearchReindexCli.PostTask(ctx, headers, jobID, taskToCreate)
	if err != nil {
		return err
	}
	postTaskRespETag := respHeaders.ETag

	logData = log.Data{
		"Task":                                 searchReindexTask,
		"ETag returned in patch job response:": patchJobRespETag,
		"ETag returned in post task response:": postTaskRespETag,
	}

	log.Info(ctx, "event successfully handled", logData)
	return nil
}

func (h *ReindexRequestedHandler) GetPayload(ctx context.Context, noOfDocument int, m map[string]string) ([]byte, error) {
	var task string
	if len(m) == 0 {
		return nil, errors.New("no tasks exists")
	}

	for key, value := range m {
		if key != zebedeeTaskName {
			log.Info(ctx, "invalid taskname", log.Data{
				"taskvalue": task,
			})
			return nil, errors.New("invalid taskname")
		}
		task = value
	}

	payload := []byte(task)
	payload = append(payload, byte(noOfDocument))
	log.Info(ctx, "endpoint request params", log.Data{
		"payload": string(payload),
	})

	return payload, nil
}
