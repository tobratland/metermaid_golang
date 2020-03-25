package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/faker"
	"github.com/tobra/metermaid/models"
	"github.com/tobra/metermaid/usecase"

	"github.com/tobra/metermaid/repository/mocks"
)

func TestStoreTimeseriesValues(t *testing.T) {
	mockTimeseriesRepo := new(mocks.TimeseriesRepository)
	var mockTimeseries models.TimeSeries
	err := faker.FakeData(&mockTimeseries)

	assert.NoError(t, err)

	tempMockTimeseries := mockTimeseries
	tempMockTimeseries.Id = "123test"
	mockTimeseriesRepo.On("Store", &tempMockTimeseries).Return(mockTimeseries.Id, nil)
	defer mockTimeseriesRepo.AssertExpectations(t)

	u := usecase.NewTimeseriesUsecase(mockTimeseriesRepo)

	err = u.Store(&tempMockTimeseries)

	assert.NoError(t, err)
}
