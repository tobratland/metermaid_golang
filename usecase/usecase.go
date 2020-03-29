package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tobra/metermaid/models"
	"github.com/tobra/metermaid/repository"
)

type TimeseriesUsecase interface {
	Store(*models.TimeSeries) error
	GetAllTimeseriesFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error)
}

type timeseriesUsecase struct {
	timeseriesRepo repository.TimeSeriesRepository
}

func NewTimeseriesUsecase(t repository.TimeSeriesRepository) TimeseriesUsecase {
	return &timeseriesUsecase{t}
}

func (t *timeseriesUsecase) Store(ts *models.TimeSeries) error {
	ts.Id = uuid.New().String()
	err := t.timeseriesRepo.StoreValues(ts)
	if err != nil {
		return err
	}

	err = t.timeseriesRepo.StoreData(ts)
	if err != nil {
		return err
	}

	return nil
}

func (t *timeseriesUsecase) GetAllTimeseriesFromTimeToTime(from time.Time, to time.Time) (returnedTs []models.TimeSeries, err error) {

	tss, err := t.timeseriesRepo.GetAllDataFromTimeToTime(from, to)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(tss); i++ {
		ts := tss[i]
		tsWithValues, err := t.timeseriesRepo.GetValuesByTimeseries(&ts)
		if err != nil {
			fmt.Println(err)
		}
		tss[i] = *tsWithValues
	}

	return tss, nil
}

func ParseStringToTime(timeString string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println(err)
	}
	return t
}
