package mocks

import (
	mock "github.com/stretchr/testify/mock"
	"github.com/tobra/metermaid/models"
)

type TimeseriesRepository struct {
	mock.Mock
}

func (_m *TimeseriesRepository) Store(t *models.TimeSeries) error {
	ret := _m.Called(t)

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.TimeSeries) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(1)
	}

	return r1
}
