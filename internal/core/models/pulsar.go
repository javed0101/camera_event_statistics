package models

type PulsarEvent struct {
	EventID   *string   `json:"eventID"`
	CamersID  *string   `json:"cameraID"`
	TimeStamp *string   `json:"timestamp"`
	Location  *Location `json:"location"`
	EventType *string   `json:"eventType"`
	MetaData  *MetaData `json:"metaData"`
}

type Location struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type MetaData struct {
	ObjectID        *string  `json:"objectID"`
	ConfidenceScore *float64 `json:"confidenceScore"`
}
