package models

// ContentUpdate provides an avro structure for a Content Update event
type ContentUpdate struct {
	URI          string `avro:"uri"`
	DataType     string `avro:"data_type"`
	CollectionID string `avro:"collection_id"`
	TraceID      string `avro:"trace_id"`
}
