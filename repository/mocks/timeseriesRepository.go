package mocks

import (
	"time"

	mock "github.com/stretchr/testify/mock"
	"github.com/tobra/metermaid/models"
)

type TimeseriesRepository struct {
	mock.Mock
}

func (_m *TimeseriesRepository) StoreValues(t *models.TimeSeries) error {
	//TODO fix the test
	ret := _m.Called(t)

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.TimeSeries) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}
func (_m *TimeseriesRepository) StoreData(t *models.TimeSeries) error {
	ret := _m.Called(t)

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.TimeSeries) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}

func (_m *TimeseriesRepository) GetAllDataFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error) {
	ret := _m.Called(from, to)

	var r0 []models.TimeSeries

	if rf, ok := ret.Get(0).(func(time.Time, time.Time) []models.TimeSeries); ok {
		r0 = rf(from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.TimeSeries)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time, time.Time) error); ok {
		r1 = rf(from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *TimeseriesRepository) GetValuesByTimeseries(t *models.TimeSeries) (*models.TimeSeries, error) {
	ret := _m.Called(t)

	var r0 *models.TimeSeries

	if rf, ok := ret.Get(0).(func(*models.TimeSeries) *models.TimeSeries); ok {
		r0 = rf(t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.TimeSeries)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.TimeSeries) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1

}

/*
StoreData(t *models.TimeSeries) error
	GetAllDataFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error)
	GetValuesByTimeseries(t *models.TimeSeries) (*models.TimeSeries, error)
*/
