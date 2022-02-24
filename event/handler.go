package event

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// TODO: remove hello called example handler
// HelloCalledHandler ...
type HelloCalledHandler struct {
}

// Handle takes a single event.
func (h *HelloCalledHandler) Handle(ctx context.Context, cfg *config.Config, event *HelloCalled) (err error) {
	logData := log.Data{
		"event": event,
	}
	log.Info(ctx, "event handler called", logData)

	greeting := fmt.Sprintf("Hello, %s!", event.RecipientName)
	err = ioutil.WriteFile(cfg.OutputFilePath, []byte(greeting), 0644)
	if err != nil {
		return err
	}

	logData["greeting"] = greeting
	log.Info(ctx, "hello world example handler called successfully", logData)
	log.Info(ctx, "event successfully handled", logData)

	return nil
}
