package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tobra/metermaid/delivery/grpc/timeseries_grpc"
	"github.com/tobra/metermaid/models"
	_usecase "github.com/tobra/metermaid/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewTimeSeriesServerGrpc(gserver *grpc.Server, TimeseriesUsecase _usecase.TimeseriesUsecase) {

	timeseriesServer := &server{
		usecase: TimeseriesUsecase,
	}

	timeseries_grpc.RegisterTimeseriesHandlerServer(gserver, timeseriesServer)
	reflection.Register(gserver)
}

type server struct {
	usecase _usecase.TimeseriesUsecase
}

func (s *server) transformTimeseriesData(ts *timeseries_grpc.Timeseries) *models.TimeSeries {
	m := make(map[time.Time]float64)
	for k, v := range ts.Values {
		m[ParseStringToTime(k)] = v
	}

	res := &models.TimeSeries{
		MeterId:    ts.MeterId,
		CustomerId: ts.CustomerId,
		Resolution: ts.Resolution,
		From:       ParseStringToTime(ts.From),
		To:         ParseStringToTime(ts.To),
		Values:     m,
	}

	return res
}

func (s *server) transformTimeseriesRpc(ts *models.TimeSeries) *timeseries_grpc.Timeseries {
	if ts == nil {
		return nil
	}
	m := make(map[string]float64)
	for k, v := range ts.Values {
		m[ParseTimeToString(k)] = v
	}
	res := &timeseries_grpc.Timeseries{
		Id:         ts.Id,
		MeterId:    ts.MeterId,
		CustomerId: ts.CustomerId,
		Resolution: ts.Resolution,
		From:       ParseTimeToString(ts.From),
		To:         ParseTimeToString(ts.To),
		Values:     m,
	}
	return res
}

func (s *server) Store(ctx context.Context, t *timeseries_grpc.Timeseries) (*timeseries_grpc.Timeseries, error) {
	ts := s.transformTimeseriesData(t)
	err := s.usecase.Store(ts)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return t, nil
}

func (s *server) GetAllFromTimeToTime(ctx context.Context, r *timeseries_grpc.GetAllFromTimeToTimeRequest) (*timeseries_grpc.GetAllFromTimeToTimeResponse, error) {
	tss, err := s.usecase.GetAllTimeseriesFromTimeToTime(ParseStringToTime(r.From), ParseStringToTime(r.To))
	res := make([]*timeseries_grpc.Timeseries, len(tss))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	for i, t := range tss {
		ts := s.transformTimeseriesRpc(&t)
		res[i] = ts
	}
	response := &timeseries_grpc.GetAllFromTimeToTimeResponse{
		TimeseriesList: res,
	}
	return response, nil
}

func ParseStringToTime(timeString string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func ParseTimeToString(t time.Time) string {
	layout := "2006-01-02T15:04:05Z"
	ts := t.Format(layout)
	return ts
}
