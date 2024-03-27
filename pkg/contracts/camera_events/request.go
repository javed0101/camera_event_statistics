package cameraevents

type CameraEventRequest struct {
	CameraID  *string `json:"cameraId" required:"true" queryparam:"cameraID"`
	EventType *string `json:"eventType" required:"true" queryparam:"eventType"`
}
