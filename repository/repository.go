package repository

import (
	"time"

	"github.com/tobra/metermaid/models"
)

type TimeSeriesRepository interface {
	/* Store(t *models.TimeSeries) (string, error) */
	StoreValues(t *models.TimeSeries) error
	StoreData(t *models.TimeSeries) error
	GetAllTimeseriesFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error)
}
