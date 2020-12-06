package workr

import "encoding/json"

// QueueItem represent an item in the queue
type QueueItem struct {
	JobName string `json:"n"`
	Data    []byte `json:"d"`
}

// Load load data bytes to *QueueItem itm
func (itm *QueueItem) Load(data []byte) error {
	return json.Unmarshal(data, itm)
}

// Bytes return json bytes for queueItem
func (itm *QueueItem) Bytes() ([]byte, error) {
	return json.Marshal(itm)
}
