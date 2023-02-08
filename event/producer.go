package event

import (
	"context"

	dpkafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out ./mock/producer.go -pkg mock . Marshaller

// Marshaller defines a type for marshalling the requested object into the required format.
type Marshaller interface {
	Marshal(s interface{}) ([]byte, error)
}

// ContentUpdatedProducer produces kafka messages for instances which have been successfully processed.
type ContentUpdatedProducer struct {
	Marshaller Marshaller
	Producer   dpkafka.IProducer
}

// ContentUpdate produce a kafka message for an instance which has been successfully processed.
func (p ContentUpdatedProducer) ContentUpdate(ctx context.Context, cfg *config.Config, event models.ContentUpdated) error {
	if cfg.EnablePublishContentUpdatedTopic {
		log.Info(ctx, "EnablePublishContentUpdatedTopic Flag is enabled")
		eventBytes, err := p.Marshaller.Marshal(event)
		if err != nil {
			log.Fatal(ctx, "failed to marshal event", err)
			return err
		}

		p.Producer.Channels().Output <- eventBytes
		log.Info(ctx, "event produced successfully", log.Data{"event": event, "package": "event.ContentUpdate"})
	} else {
		log.Info(ctx, "EnablePublishContentUpdatedTopic Flag is disabled")
	}
	return nil
}
