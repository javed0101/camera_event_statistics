package cameraevents

type CameraEventResponse struct {
	RequestID       *string    `json:"requestID"`
	EventCount      *int       `json:"eventCount"`
	TotalEventCount *int       `json:"totalEventCount"`
	TimeStamp       *TimeStamp `json:"timestamp"`
}

type TimeStamp struct {
	FirstEvent *string `json:"firstEvent"`
	LastEvent  *string `json:"lastEvent"`
}
