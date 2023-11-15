package handler

import (
	"context"
	"strings"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/log.go/v2/log"
)

var (
	defaultGoRoutineCount    = 10
	maxConcurrentExtractions = 20
	DefaultPaginationLimit   = 500
)

const (
	TotalTasksCount = 2
)

type DatasetEditionMetadata struct {
	id        string
	editionID string
	version   string
}

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeCli                clients.ZebedeeClient
	DatasetAPICli             clients.DatasetAPIClient
	ContentUpdatedProducer    event.ContentUpdatedProducer
	ReindexTaskCountsProducer event.ReindexTaskCountsProducer
	Config                    *config.Config
}

type taskDetails struct {
	Name  string
	count int
}

// Handle takes a single event.
// TODO - any error which occurs in the ReindexRequestedHandler.Handle function are not returned and only logged
// We are only logging as we are not handling errors via an error topic yet
func (h *ReindexRequestedHandler) Handle(ctx context.Context, reindexReqEvent *models.ReindexRequested) {
	logData := log.Data{
		"event": reindexReqEvent,
	}
	log.Info(ctx, "reindex requested event handler called", logData)

	taskCounter := 1
	tasks := make(chan taskDetails)
	var extractionCompleted bool
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		h.getAndSendZebedeeDocsURL(ctx, h.Config, h.ZebedeeCli, reindexReqEvent, tasks)
	}()
	go func() {
		defer wg.Done()
		h.getAndSendDatasetURLs(ctx, h.Config, h.DatasetAPICli, reindexReqEvent, tasks)
	}()
	go func() {
		defer wg.Done()
		for task := range tasks {
			if taskCounter == TotalTasksCount {
				extractionCompleted = true
			}
			err := h.ReindexTaskCountsProducer.TaskCounts(ctx, h.Config, models.ReindexTaskCounts{
				JobID:               reindexReqEvent.JobID,
				Task:                task.Name,
				Count:               int32(task.count),
				ExtractionCompleted: extractionCompleted,
			})
			if err != nil {
				log.Error(ctx, "failed to publish task counts to reindex task counts topic due to", err, log.Data{"request_params": nil})
				return
			}
			taskCounter++
			if taskCounter > TotalTasksCount {
				break
			}
		}
	}()
	wg.Wait()

	log.Info(ctx, "event successfully handled", logData)
}

func (h *ReindexRequestedHandler) getAndSendZebedeeDocsURL(ctx context.Context, cfg *config.Config, zebedeeCli clients.ZebedeeClient, reindexReqEvent *models.ReindexRequested, task chan taskDetails) {
	log.Info(ctx, "extract and publish zebedee docs url")
	publishedIndex, err := zebedeeCli.GetPublishedIndex(ctx, nil)
	if err != nil {
		// cfg.ZebedeeClientTimeout may need to be increased
		// Note that this call is to retrieve all published files on disk, this operation in Zebedee takes a long time to process and respond to caller (search data finder)
		// Problem is that as more and more data/content is published, more data is stored on disk in Zebedee and hence the time to process will increase
		log.Error(ctx, "failed to get published index from zebedee", err, log.Data{"request_params": nil})
		return
	}

	urlList := make([]string, 10)
	publishedItems := publishedIndex.Items
	totalZebedeeDocs := len(publishedItems)

	logData := log.Data{
		"totaldocs": totalZebedeeDocs,
	}

	log.Info(ctx, "publish to content updated topic", logData)

	// it takes more than 10 mins to retrieve all document urls from zebedee
	// use ZEBEDEE_REQUEST_LIMIT to limit this for performance / testing / debugging
	limit := 0

	if cfg.ZebedeeRequestLimit == 0 || cfg.ZebedeeRequestLimit > totalZebedeeDocs {
		limit = totalZebedeeDocs
	} else {
		limit = cfg.ZebedeeRequestLimit
	}

	for i := 0; i < limit; i++ {
		err := h.ContentUpdatedProducer.ContentUpdate(ctx, cfg, models.ContentUpdated{
			URI:         publishedIndex.Items[i].URI,
			JobID:       reindexReqEvent.JobID,
			TraceID:     reindexReqEvent.TraceID,
			SearchIndex: reindexReqEvent.SearchIndex,
		})
		if err != nil {
			log.Error(ctx, "failed to publish zebedee doc to content updated topic", err, log.Data{"request_params": nil})
			return
		}
	}

	log.Info(ctx, "first 10 Zebedee docs URLs retrieved", log.Data{"first URLs": urlList})
	log.Info(ctx, "total number of zebedee docs retrieved", log.Data{"total_documents": totalZebedeeDocs})
	taskNames := strings.Split(cfg.TaskNameValues, ",")
	task <- taskDetails{
		Name:  extractTaskName(taskNames, "zebedee"),
		count: totalZebedeeDocs,
	}
}

func (h *ReindexRequestedHandler) getAndSendDatasetURLs(ctx context.Context, cfg *config.Config, datasetAPICli clients.DatasetAPIClient, reindexReqEvent *models.ReindexRequested, task chan taskDetails) {
	log.Info(ctx, "extract and send dataset urls")
	var wgDataset sync.WaitGroup
	wgDataset.Add(4)
	datasetChan := h.extractDatasets(ctx, &wgDataset, datasetAPICli, cfg.ServiceAuthToken)
	editionChan := h.retrieveDatasetEditions(ctx, &wgDataset, datasetAPICli, datasetChan, cfg.ServiceAuthToken)
	datasetURLChan := h.getAndSendDatasetURLsFromLatestMetadata(ctx, &wgDataset, datasetAPICli, editionChan, cfg.ServiceAuthToken)
	urlCount := h.sendExtractedDatasetURLs(ctx, &wgDataset, cfg, datasetURLChan, reindexReqEvent)
	wgDataset.Wait() // wait for the other go-routines to complete which extracts the dataset urls
	log.Info(ctx, "successfully extracted all datasets")
	taskNames := strings.Split(cfg.TaskNameValues, ",")
	task <- taskDetails{
		Name:  extractTaskName(taskNames, "dataset"),
		count: urlCount,
	}
}

func (h *ReindexRequestedHandler) extractDatasets(ctx context.Context, wgDataset *sync.WaitGroup, datasetAPIClient clients.DatasetAPIClient, serviceAuthToken string) chan dataset.Dataset {
	log.Info(ctx, "extract datasets")
	datasetChan := make(chan dataset.Dataset)
	go func() {
		defer close(datasetChan)
		defer wgDataset.Done()
		var offset = 0
		var totalDocs = 0
		for {
			list, err := datasetAPIClient.GetDatasets(ctx, "", serviceAuthToken, "", &dataset.QueryParams{
				Offset: offset,
				Limit:  DefaultPaginationLimit,
			})
			if err != nil {
				log.Error(ctx, "failed to get dataset clients: %v", err)
				break
			}

			listLength := len(list.Items)

			log.Info(ctx, "datasets list length", log.Data{"url list length": listLength})

			if listLength == 0 {
				break
			}

			for i := 0; i < listLength; i++ {
				datasetChan <- list.Items[i]
			}

			offset += DefaultPaginationLimit
			totalDocs += listLength
		}
		log.Info(ctx, "total number of dataset docs retrieved", log.Data{"total_documents": totalDocs})
	}()
	return datasetChan
}

func (h *ReindexRequestedHandler) retrieveDatasetEditions(ctx context.Context, wgDataset *sync.WaitGroup, datasetAPIClient clients.DatasetAPIClient, datasetChan chan dataset.Dataset, serviceAuthToken string) chan DatasetEditionMetadata {
	log.Info(ctx, "retrieve dataset editions")
	editionMetadataChan := make(chan DatasetEditionMetadata)
	go func() {
		defer close(editionMetadataChan)
		defer wgDataset.Done()
		var wg sync.WaitGroup
		for i := 0; i < defaultGoRoutineCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for dataset := range datasetChan {
					if dataset.Current == nil {
						continue
					}
					editions, err := datasetAPIClient.GetFullEditionsDetails(ctx, "", serviceAuthToken, dataset.CollectionID, dataset.Current.ID)
					if err != nil {
						logData := log.Data{
							"dataset_id":    dataset.Current.ID,
							"collection_id": dataset.CollectionID,
						}
						log.Error(ctx, "error retrieving editions for dataset", err, logData)
					}
					for i := 0; i < len(editions); i++ {
						if editions[i].ID == "" || editions[i].Current.Links.LatestVersion.ID == "" {
							continue
						}
						editionMetadataChan <- DatasetEditionMetadata{
							id:        dataset.Current.ID,
							editionID: editions[i].Current.Edition,
							version:   editions[i].Current.Links.LatestVersion.ID,
						}
					}
				}
			}()
		}
		wg.Wait()
		log.Info(ctx, "successfully retrieved dataset editions")
	}()
	return editionMetadataChan
}

func (h *ReindexRequestedHandler) getAndSendDatasetURLsFromLatestMetadata(ctx context.Context, wgDataset *sync.WaitGroup, datasetAPIClient clients.DatasetAPIClient, editionMetadata chan DatasetEditionMetadata, serviceAuthToken string) chan string {
	log.Info(ctx, "extract and send dataset url from the latest metadata")
	datasetURLChan := make(chan string)
	go func() {
		defer close(datasetURLChan)
		defer wgDataset.Done()
		var wg sync.WaitGroup
		for i := 0; i < defaultGoRoutineCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for edMetadata := range editionMetadata {
					log.Info(ctx, "received edition metadata", log.Data{"metadata": edMetadata})
					metadata, err := datasetAPIClient.GetVersionMetadata(ctx, "", serviceAuthToken, "", edMetadata.id, edMetadata.editionID, edMetadata.version)
					if err != nil {
						continue
					}
					url := metadata.DatasetLinks.LatestVersion.URL
					datasetURLChan <- url
				}
			}()
		}
		log.Info(ctx, "successfully handled metadata")
		wg.Wait()
	}()
	return datasetURLChan
}

func (h *ReindexRequestedHandler) sendExtractedDatasetURLs(ctx context.Context, wgDataset *sync.WaitGroup, cfg *config.Config, datasetURLChan chan string, reindexReqEvent *models.ReindexRequested) int {
	var urlListCount int
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer wgDataset.Done()
		for datasetURL := range datasetURLChan {
			log.Info(ctx, "send extracted dataset urls")
			err := h.ContentUpdatedProducer.ContentUpdate(ctx, cfg, models.ContentUpdated{
				URI:         datasetURL,
				DataType:    DatasetDataType,
				JobID:       reindexReqEvent.JobID,
				TraceID:     reindexReqEvent.TraceID,
				SearchIndex: reindexReqEvent.SearchIndex,
			})
			if err != nil {
				log.Error(ctx, "failed to publish datasets to content update topic", err)
				return
			}
			urlListCount++
		}
	}()
	wg.Wait()
	return urlListCount
}

func extractTaskName(taskNames []string, matchingName string) string {
	for _, taskName := range taskNames {
		if strings.Contains(taskName, matchingName) {
			return taskName
		}
	}
	return ""
}
