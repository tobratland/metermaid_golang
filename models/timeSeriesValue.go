package models

import "time"

type TimeSeriesValue struct {
	Id           string
	TimeSeriesId string
	MeterId      string
	CustomerId   string
	Hour         time.Time
	Value        float32
}

func NewTimeSeriesValue(id string, timeSeriesId string, meterId string, customerId string, hour time.Time, value float32) *TimeSeriesValue {
	return &TimeSeriesValue{id, timeSeriesId, meterId, customerId, hour, value}
}
