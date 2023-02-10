package handler

import (
	"context"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/log.go/v2/log"
)

var (
	maxConcurrentExtractions = 20
	DefaultPaginationLimit   = 500
)

type DatasetEditionMetadata struct {
	id        string
	editionID string
	version   string
}

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	ZebedeeCli    clients.ZebedeeClient
	DatasetAPICli clients.DatasetAPIClient
	Producer      event.ContentUpdatedProducer
	Config        *config.Config
}

// Handle takes a single event.
// TODO - any error which occurs in the ReindexRequestedHandler.Handle function are not returned and only logged
// We are only logging as we are not handling errors via an error topic yet
func (h *ReindexRequestedHandler) Handle(ctx context.Context, reindexReqEvent *models.ReindexRequested) {
	logData := log.Data{
		"event": reindexReqEvent,
	}
	log.Info(ctx, "reindex requested event handler called", logData)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		h.getAndSendZebedeeDocsURL(ctx, h.Config, h.ZebedeeCli, reindexReqEvent)
	}()
	go func() {
		defer wg.Done()
		h.getAndSendDatasetURLs(ctx, h.Config, h.DatasetAPICli, reindexReqEvent)
	}()
	wg.Wait()

	log.Info(ctx, "event successfully handled", logData)
}

func (h *ReindexRequestedHandler) getAndSendZebedeeDocsURL(ctx context.Context, cfg *config.Config, zebedeeCli clients.ZebedeeClient, reindexReqEvent *models.ReindexRequested) {
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

	// only the first 10 docs are retrieved for testing and performance purposes
	// it takes more than 10 mins to retrieve all document urls from zebedee
	// TODO: remove (i < 10) condition when this app has been completely implemented
	for i := 0; (i < 10) && (i < totalZebedeeDocs); i++ {
		err := h.Producer.ContentUpdate(ctx, cfg, models.ContentUpdated{
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
}

func (h *ReindexRequestedHandler) getAndSendDatasetURLs(ctx context.Context, cfg *config.Config, datasetAPICli clients.DatasetAPIClient, reindexReqEvent *models.ReindexRequested) {
	var wgDataset sync.WaitGroup
	wgDataset.Add(4)
	datasetChan := h.extractDatasets(ctx, &wgDataset, datasetAPICli, cfg.ServiceAuthToken)
	editionChan := h.retrieveDatasetEditions(ctx, &wgDataset, datasetAPICli, datasetChan, cfg.ServiceAuthToken)
	datasetURLChan := h.getAndSendDatasetURLsFromLatestMetadata(ctx, &wgDataset, datasetAPICli, editionChan, cfg.ServiceAuthToken)
	h.logExtractedDatasetURLs(ctx, &wgDataset, cfg, datasetURLChan, reindexReqEvent)
	// TODO - logExtractedDatasetURLs is temporary and should be replaced in the future
	wgDataset.Wait() // wait for the other go-routines to complete which extracts the dataset urls
}

func (h *ReindexRequestedHandler) extractDatasets(ctx context.Context, wgDataset *sync.WaitGroup, datasetAPIClient clients.DatasetAPIClient, serviceAuthToken string) chan dataset.Dataset {
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
	editionMetadataChan := make(chan DatasetEditionMetadata)
	go func() {
		defer close(editionMetadataChan)
		defer wgDataset.Done()
		noOfConcurrentExtractions := getNoOfConcurrentExtractions(len(datasetChan), maxConcurrentExtractions)
		var wg sync.WaitGroup
		for i := 0; i < noOfConcurrentExtractions; i++ {
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
	}()
	return editionMetadataChan
}

func (h *ReindexRequestedHandler) getAndSendDatasetURLsFromLatestMetadata(ctx context.Context, wgDataset *sync.WaitGroup, datasetAPIClient clients.DatasetAPIClient, editionMetadata chan DatasetEditionMetadata, serviceAuthToken string) chan string {
	datasetURLChan := make(chan string)
	go func() {
		defer close(datasetURLChan)
		defer wgDataset.Done()
		var wg sync.WaitGroup
		noOfConcurrentExtractions := getNoOfConcurrentExtractions(len(editionMetadata), maxConcurrentExtractions)
		for i := 0; i < noOfConcurrentExtractions; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for edMetadata := range editionMetadata {
					metadata, err := datasetAPIClient.GetVersionMetadata(ctx, "", serviceAuthToken, "", edMetadata.id, edMetadata.editionID, edMetadata.version)
					if err != nil {
						continue
					}
					url := metadata.DatasetLinks.LatestVersion.URL
					datasetURLChan <- url
				}
			}()
		}
		wg.Wait()
	}()
	return datasetURLChan
}

// TODO - logExtractedDatasetURLs is temporary.
// The dataset url should be sent to the content-updated topic here in the future.
// But for the time being, we are going to extract the urls and print them
func (h *ReindexRequestedHandler) logExtractedDatasetURLs(ctx context.Context, wgDataset *sync.WaitGroup, cfg *config.Config, datasetURLChan chan string, reindexReqEvent *models.ReindexRequested) {
	urlList := make([]string, 0)
	go func() {
		defer wgDataset.Done()
		for datasetURL := range datasetURLChan {
			err := h.Producer.ContentUpdate(ctx, cfg, models.ContentUpdated{
				URI:         datasetURL,
				JobID:       reindexReqEvent.JobID,
				TraceID:     reindexReqEvent.TraceID,
				SearchIndex: reindexReqEvent.SearchIndex,
			})
			if err != nil {
				log.Error(ctx, "failed to publish datasets to content update topic", err)
				return
			}
		}
	}()

	log.Info(ctx, "dataset docs URLs retrieved", log.Data{"urls": urlList})
}

func getNoOfConcurrentExtractions(chanLength, max int) int {
	if chanLength < max {
		return chanLength
	}
	return max
}
