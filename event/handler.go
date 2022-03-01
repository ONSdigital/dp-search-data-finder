package event

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// ReindexRequestedHandler ...
type ReindexRequestedHandler struct {
}

// Handle takes a single event.
func (h *ReindexRequestedHandler) Handle(ctx context.Context, cfg *config.Config, event *ReindexRequested) (err error) {
	logData := log.Data{
		"event": event,
	}
	log.Info(ctx, "event handler called", logData)

	greeting := fmt.Sprintf("Hello there! Job id is %s", event.JobID)
	err = ioutil.WriteFile(cfg.OutputFilePath, []byte(greeting), 0644)
	if err != nil {
		return err
	}

	logData["greeting"] = greeting
	log.Info(ctx, "reindex requested example handler called successfully", logData)
	log.Info(ctx, "event successfully handled", logData)

	return nil
}
