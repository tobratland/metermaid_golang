package usecase

import (
	"github.com/google/uuid"
	"github.com/tobra/metermaid/models"
	"github.com/tobra/metermaid/repository"
)

type TimeseriesUsecase interface {
	Store(*models.TimeSeries) error
}

type timeseriesUsecase struct {
	timeseriesRepo repository.TimeSeriesRepository
}

func NewTimeseriesUsecase(t repository.TimeSeriesRepository) TimeseriesUsecase {
	return &timeseriesUsecase{t}
}

func (t *timeseriesUsecase) Store(m *models.TimeSeries) error {
	m.Id = uuid.New().String()
	err := t.timeseriesRepo.Store(m)
	if err != nil {
		return err
	}

	return nil
}
