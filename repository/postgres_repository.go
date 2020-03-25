package repository

import (
	"database/sql"
	"fmt"
	"log"

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

/* func (p *postgresTimeseriesRepository) Store(t *models.TimeSeries) (string, error) {
	query := `INSERT INTO timeseries(id, meterid, customerid, resolution, fromtime, totime) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	stmt, err := p.Conn.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return "", models.INTERNAL_SERVER_ERROR
	}
	var id string
	err = stmt.QueryRow(t.Id, t.MeterId, t.CustomerId, t.Resolution, t.From, t.To).Scan(&id)
	if err != nil {
		log.Fatal(err)
		return "", models.INTERNAL_SERVER_ERROR
	}
	return id, nil

} */

func (p *postgresTimeseriesRepository) Store(t *models.TimeSeries) error {
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
	fmt.Println("successfully stored")
	return nil
}

/* Id           string
TimeSeriesId string
MeterId      string
CustomerId   string
Hour         time.Time
Value        float32 */
//insert into meterdata(id, meterid, userid, resolution, fromtime, totime) values(?,?,?,?,?,?)
