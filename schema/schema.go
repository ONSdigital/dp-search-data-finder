package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var contentUpdatedEvent = `{
  "type": "record",
  "name": "content-updated",
  "fields": [
    {"name": "uri", "type": "string", "default": ""},
    {"name": "data_type", "type": "string", "default": ""},
    {"name": "collection_id", "type": "string", "default": ""},
    {"name": "job_id", "type": "string", "default": ""},
    {"name": "trace_id", "type": "string", "default": ""},
    {"name": "search_index", "type": "string", "default": ""}
  ]
}`

// ContentUpdatedEvent is the Avro schema for Content Updated messages.
var ContentUpdatedEvent = &avro.Schema{
	Definition: contentUpdatedEvent,
}

var reindexRequestedEvent = `{
	"type": "record",
	"name": "reindex-requested",
	"fields": [
		{"name": "job_id", "type": "string", "default": ""},
		{"name": "search_index", "type": "string", "default": ""},
		{"name": "trace_id", "type": "string", "default": ""}
	]
}`

// ReindexRequestedEvent is the Avro schema for Reindex Requested messages.
var ReindexRequestedEvent = &avro.Schema{
	Definition: reindexRequestedEvent,
}
