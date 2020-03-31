package repository

import (
	"time"

	"github.com/tobra/metermaid/models"
)

type TimeSeriesRepository interface {
	/* Store(t *models.TimeSeries) (string, error) */
	StoreValues(t *models.TimeSeries) error
	StoreData(t *models.TimeSeries) error
	GetAllDataFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error)
	GetSumFromTimeToTimeByCustomerId(from time.Time, to time.Time, customerId string) (float64, error)
	GetValuesByTimeseries(t *models.TimeSeries) (*models.TimeSeries, error)
}
