package helper

import "time"

func StringPointer(s string) *string {
	return &s
}

func IntPointer(i int) *int {
	return &i
}

func FloatPointer(f float64) *float64 {
	return &f
}

func MinTimeStamp(timestamp1, timestamp2 string) string {
	time1, _ := time.Parse(time.RFC3339, timestamp1)
	time2, _ := time.Parse(time.RFC3339, timestamp2)

	if time1.Before(time2) {
		return time1.Format(time.RFC3339)
	}
	return time2.Format(time.RFC3339)
}
