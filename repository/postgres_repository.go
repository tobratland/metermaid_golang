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
		fmt.Println("conn.begin @ storedata")
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
	fmt.Println(len(t.Values))
	for key := range t.Values {
		_, err = stmt.Exec(uuid.New().String(), t.Id, t.MeterId, t.CustomerId, key, t.Values[key])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(key, " ", t.Values[key])
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

func (p *postgresTimeseriesRepository) GetAllTimeseriesFromTimeToTime(from time.Time, to time.Time) (timeseries []models.TimeSeries, err error) {
	/* txn, err := p.Conn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer p.Conn.Close()
	statement := "SELECT * FROM metervalues WHERE hour >=$1 AND hour <= $2"
	rows, err := txn.Query(statement, from, to)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var timeseries models.TimeSeries

		if err := rows.Scan(&timeseries.Id, &timeseries.); err != nil {
			log.Fatal(err)
		}
	} */

	return
}

/* Id           string
TimeSeriesId string
MeterId      string
CustomerId   string
Hour         time.Time
Value        float32 */
//insert into meterdata(id, meterid, userid, resolution, fromtime, totime) values(?,?,?,?,?,?)
