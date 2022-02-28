package event

// ContentUpdated provides an avro structure for a Content Updated event
// TODO: Add the correct fields to the ContentUpdated struct
type ContentUpdated struct {
	RecipientName string `avro:"recipient_name"`
}

// ReindexRequested provides an avro structure for a Reindex Requested event
type ReindexRequested struct {
	JobID       string `avro:"job_id"`
	SearchIndex string `avro:"search_index"`
	TraceID     string `avro:"trace_id"`
}
