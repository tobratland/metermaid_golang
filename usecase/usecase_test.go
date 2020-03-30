package usecase_test

import (
	"testing"
	"time"

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
	mockTimeseriesRepo.On("StoreValues", &tempMockTimeseries).Return(nil)
	mockTimeseriesRepo.On("StoreData", &tempMockTimeseries).Return(nil)
	defer mockTimeseriesRepo.AssertExpectations(t)

	u := usecase.NewTimeseriesUsecase(mockTimeseriesRepo)

	err = u.Store(&tempMockTimeseries)

	assert.NoError(t, err)
}

func TestGetAllTimeseriesFromTimeToTime(t *testing.T) {
	mockTimeseriesRepo := new(mocks.TimeseriesRepository)
	var mockTimeseries models.TimeSeries
	err := faker.FakeData(&mockTimeseries)
	assert.NoError(t, err)
	var mockTimeseriesList []models.TimeSeries
	mockTimeseriesList = append(mockTimeseriesList, mockTimeseries)

	from, err := time.Parse("2006-01-02T15:04:05Z", "1988-11-25T00:00:00Z")
	assert.NoError(t, err)

	to, err := time.Parse("2006-01-02T15:04:05Z", "1988-11-25T23:00:00Z")
	assert.NoError(t, err)

	mockTimeseriesRepo.On("GetAllDataFromTimeToTime", from, to).Return(mockTimeseriesList, err)
	mockTimeseriesRepo.On("GetValuesByTimeseries", &mockTimeseries).Return(&mockTimeseries, err)

	u := usecase.NewTimeseriesUsecase(mockTimeseriesRepo)

	list, err := u.GetAllTimeseriesFromTimeToTime(from, to)
	expected := mockTimeseriesList

	assert.Equal(t, list, expected)
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, len(mockTimeseriesList))

	mockTimeseriesRepo.AssertCalled(t, "GetAllDataFromTimeToTime", from, to)
	mockTimeseriesRepo.AssertCalled(t, "GetValuesByTimeseries", &mockTimeseries)

}
