package event

import (
	"context"
	"fmt"

	dpkafka "github.com/ONSdigital/dp-kafka/v2"
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
func (p ContentUpdatedProducer) ContentUpdate(ctx context.Context, event models.ContentUpdated) error {
	bytes, err := p.Marshaller.Marshal(event)
	if err != nil {
		log.Fatal(ctx, "Marshaller.Marshal", err)
		return fmt.Errorf(fmt.Sprintf("Marshaller.Marshal returned an error: event=%v: %%w", event), err)
	}

	p.Producer.Channels().Output <- bytes
	log.Info(ctx, "completed successfully", log.Data{"event": event, "package": "event.ContentUpdate"})
	return nil
}
