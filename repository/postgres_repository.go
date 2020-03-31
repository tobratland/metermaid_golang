package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/tobra/metermaid/models"
)

type postgresTimeseriesRepository struct {
	Conn *sql.DB
}

func NewPostgresTimeSeriesRepository(Conn *sql.DB) TimeSeriesRepository {
	return &postgresTimeseriesRepository{Conn}
}
func (p *postgresTimeseriesRepository) StoreData(t *models.TimeSeries) error {
	txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}

	statement := "INSERT INTO meterdata (id, meterid, userid, resolution, fromtime, totime) VALUES ($1, $2, $3, $4, $5, $6) returning id"
	_, err = txn.Exec(statement, t.Id, t.MeterId, t.CustomerId, t.Resolution, t.From, t.To)
	if err != nil {
		fmt.Println("query @ storedata")
		log.Fatal(err)
		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("successfully stored data")
	return nil
}

func (p *postgresTimeseriesRepository) StoreValues(t *models.TimeSeries) error {
	txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := txn.Prepare(pq.CopyIn("metervalues", "id", "meterdataid", "meterid", "userid", "hour", "value"))
	if err != nil {
		log.Fatal(err)
	}
	for key := range t.Values {
		_, err = stmt.Exec(uuid.New().String(), t.Id, t.MeterId, t.CustomerId, key, t.Values[key])
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully stored values")
	return nil
}

func (p *postgresTimeseriesRepository) GetAllDataFromTimeToTime(from time.Time, to time.Time) (timeseries []models.TimeSeries, err error) {
	txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	statement := "SELECT * FROM meterdata WHERE fromtime >=$1 AND totime <= $2"
	rows, err := txn.Query(statement, from, to)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tempTs models.TimeSeries

		if err := rows.Scan(&tempTs.Id, &tempTs.MeterId, &tempTs.CustomerId, &tempTs.Resolution, &tempTs.From, &tempTs.To); err != nil {
			log.Fatal(err)
			return nil, err
		}
		timeseries = append(timeseries, tempTs)
	}

	return timeseries, nil
}

func (p *postgresTimeseriesRepository) GetValuesByTimeseries(t *models.TimeSeries) (*models.TimeSeries, error) {
	txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	statement := "SELECT hour, value FROM metervalues WHERE meterdataid = $1"
	rows, err := txn.Query(statement, t.Id)
	defer rows.Close()

	var hour time.Time
	var value float64
	m := make(map[time.Time]float64)
	for rows.Next() {
		rows.Scan(&hour, &value)
		m[hour] = value
	}
	t.Values = m

	return t, nil
}

func (p *postgresTimeseriesRepository) GetSumFromTimeToTimeByCustomerId(from time.Time, to time.Time, customerId string) (float64, error) {
	txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	statement := "SELECT SUM(value) FROM metervalues WHERE hour >=$1 AND hour <= $2 AND userid = $3"
	rows, err := txn.Query(statement, from, to, customerId)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer rows.Close()

	var tempsum, sum float64

	for rows.Next() {
		if err := rows.Scan(&tempsum); err != nil {
			log.Fatal(err)
			return 0, err
		}
		sum += tempsum
	}
	return sum, nil
}

func (p *postgresTimeseriesRepository) GetSumFromTimeToTimeByMeterId(from time.Time, to time.Time, meterId string) (float64, error) {
	txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	statement := "SELECT SUM(value) FROM metervalues WHERE hour >=$1 AND hour <= $2 AND meterid = $3"
	rows, err := txn.Query(statement, from, to, meterId)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer rows.Close()

	var tempsum, sum float64

	for rows.Next() {
		if err := rows.Scan(&tempsum); err != nil {
			log.Fatal(err)
			return 0, err
		}
		sum += tempsum
	}
	return sum, nil
}

/*
var u1, u2, u3, u4 string

	var ht time.Time
	var v float64
	var ts = models.TimeSeries{
		Id:         t.Id,
		MeterId:    t.MeterId,
		CustomerId: t.CustomerId,
		Resolution: t.Resolution,
		From:       t.From,
		To:         t.To,
		Values:     make(map[time.Time]float64),
	}
*/

//"id", "meterdataid", "meterid", "userid", "hour", "value"

/* Id           string
TimeSeriesId string
MeterId      string
CustomerId   string
Hour         time.Time
Value        float32 */
//insert into meterdata(id, meterid, userid, resolution, fromtime, totime) values(?,?,?,?,?,?)
func ParseStringToTime(timeString string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println(err)
	}
	return t
}
