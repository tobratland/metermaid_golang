package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tobra/metermaid/models"
	"github.com/tobra/metermaid/repository"
)

type TimeseriesUsecase interface {
	Store(*models.TimeSeries) error
	GetAllTimeseriesFromTimeToTime(from time.Time, to time.Time) ([]models.TimeSeries, error)
	GetTotalUsageForCustomerInTimePeriod(from time.Time, to time.Time, customerId string) (float64, error)
	GetTotalUsageForMeterInTimePeriod(from time.Time, to time.Time, meterId string) (float64, error)
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

func (t *timeseriesUsecase) GetTotalUsageForCustomerInTimePeriod(from time.Time, to time.Time, customerId string) (float64, error) {
	sum, err := t.timeseriesRepo.GetSumFromTimeToTimeByCustomerId(from, to, customerId)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return sum, nil
}

func (t *timeseriesUsecase) GetTotalUsageForMeterInTimePeriod(from time.Time, to time.Time, meterId string) (float64, error) {
	sum, err := t.timeseriesRepo.GetSumFromTimeToTimeByMeterId(from, to, meterId)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return sum, nil
}
