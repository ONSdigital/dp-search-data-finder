package models

// ContentUpdate provides an avro structure for a Content Update event
type ContentUpdated struct {
	URI         string `avro:"uri"`
	DataType    string `avro:"data_type"`
	JobID       string `avro:"job_id"`
	TraceID     string `avro:"trace_id"`
	SearchIndex string `avro:"search_index"`
}

// ReindexRequested provides an avro structure for a Reindex Requested event
type ReindexRequested struct {
	JobID       string `avro:"job_id"`
	SearchIndex string `avro:"search_index"`
	TraceID     string `avro:"trace_id"`
}
