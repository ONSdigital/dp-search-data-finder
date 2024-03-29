package models

// ContentUpdated provides an avro structure for a Content Updated event
type ContentUpdated struct {
	URI          string `avro:"uri"`
	DataType     string `avro:"data_type"`
	CollectionID string `avro:"collection_id"`
	JobID        string `avro:"job_id"`
	TraceID      string `avro:"trace_id"`
	SearchIndex  string `avro:"search_index"`
}

// ReindexTaskCounts provides an avro structure for a Reindex task counts event
type ReindexTaskCounts struct {
	JobID               string `avro:"job_id"`
	Task                string `avro:"task"`
	ExtractionCompleted bool   `avro:"extraction_completed"`
	Count               int32  `avro:"count"`
}

// ReindexRequested provides an avro structure for a Reindex Requested event
type ReindexRequested struct {
	JobID       string `avro:"job_id"`
	SearchIndex string `avro:"search_index"`
	TraceID     string `avro:"trace_id"`
}
