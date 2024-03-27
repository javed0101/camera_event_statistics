package models

type RedisEvent struct {
	Count     *int    `json:"count,omitempty" redis:"count,omitempty"`
	StartTime *string `json:"startTime,string,omitempty" redis:"startTime,omitempty"`
	EndTime   *string `json:"endTime,string,omitempty" redis:"endTime,omitempty"`
}
