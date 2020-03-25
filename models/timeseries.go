package models

import "time"

type TimeSeries struct {
	Id         string
	MeterId    string
	CustomerId string
	Resolution string
	From       time.Time
	To         time.Time
	Values     map[time.Time]float64
}

func NewTimeSeries(id string, meterId string, customerId string, resolution string, from time.Time, to time.Time, values map[time.Time]float64) *TimeSeries {
	return &TimeSeries{id, meterId, customerId, resolution, from, to, values}
}
