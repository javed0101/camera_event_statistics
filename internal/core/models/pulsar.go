package models

type EventInfo struct {
	LastCztsTs        *string `json:"last_czts_ts,omitempty"`
	CurrentTimestamp  *int64  `json:"currentTimestamp,omitempty"`
	CurrentThumbCount *int    `json:"currentThumbCount,omitempty"`
	SubEventType      *string `json:"subEventType,omitempty"`
}

type Event struct {
	Event     *string    `json:"event,omitempty"`
	EventType *string    `json:"eventType,omitempty"`
	EventInfo *EventInfo `json:"eventInfo,omitempty"`
	CameraID  *string    `json:"cameraid,omitempty"`
	AccountID *string    `json:"accountid,omitempty"`
	Timestamp *string    `json:"timestamp,omitempty"`
	TaskID    *string    `json:"taskid,omitempty"`
	EventID   *string    `json:"eventid,omitempty"`
}

type Info struct {
	Event *Event `json:"event,omitempty"`
}

type PulsarEvent struct {
	Info *Info `json:"info,omitempty"`
}
