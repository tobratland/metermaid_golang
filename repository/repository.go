package repository

import (
	"github.com/tobra/metermaid/models"
)

type TimeSeriesRepository interface {
	/* Store(t *models.TimeSeries) (string, error) */
	Store(t *models.TimeSeries) error
}
