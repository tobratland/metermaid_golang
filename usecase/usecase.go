package usecase

import (
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

func (t *timeseriesUsecase) GetAllTimeseriesFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error) {

	tss, err := t.timeseriesRepo.GetAllTimeseriesFromTimeToTime(from, to)
	if err != nil {
		return nil, err
	}

	return tss, nil
}
