package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
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
	for _, v := range ts.Values {

		fromParsed, err := ptypes.Timestamp(v.Hour)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		m[fromParsed] = v.Value
	}
	fromParsed, err := ptypes.Timestamp(ts.From)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	toParsed, err := ptypes.Timestamp(ts.To)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	res := &models.TimeSeries{
		MeterId:    ts.MeterId,
		CustomerId: ts.CustomerId,
		Resolution: ts.Resolution,
		From:       fromParsed,
		To:         toParsed,
		Values:     m,
	}

	return res
}

func (s *server) transformTimeseriesRpc(ts *models.TimeSeries) *timeseries_grpc.Timeseries {
	if ts == nil {
		return nil
	}

	fromParsed, err := ptypes.TimestampProto(ts.From)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	toParsed, err := ptypes.TimestampProto(ts.To)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	res := &timeseries_grpc.Timeseries{
		Id:         ts.Id,
		MeterId:    ts.MeterId,
		CustomerId: ts.CustomerId,
		Resolution: ts.Resolution,
		From:       fromParsed,
		To:         toParsed,
	}

	for k, v := range ts.Values {
		kParsed, err := ptypes.TimestampProto(k)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		res.Values = append(res.Values, &timeseries_grpc.Value{
			Hour:  kParsed,
			Value: v,
		})

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

func (s *server) GetAllFromTimeToTime(ctx context.Context, r *timeseries_grpc.GetRequest) (*timeseries_grpc.GetAllFromTimeToTimeResponse, error) {

	toParsed, err := ptypes.Timestamp(r.To)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	fromParsed, err := ptypes.Timestamp(r.From)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	tss, err := s.usecase.GetAllTimeseriesFromTimeToTime(fromParsed, toParsed)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res := make([]*timeseries_grpc.Timeseries, len(tss))

	for i, t := range tss {
		ts := s.transformTimeseriesRpc(&t)
		res[i] = ts
	}
	response := &timeseries_grpc.GetAllFromTimeToTimeResponse{
		TimeseriesList: res,
	}
	return response, nil
}

func (s *server) GetTotalUsageForCustomerInTimePeriod(ctx context.Context, r *timeseries_grpc.GetRequest) (*timeseries_grpc.SumResponse, error) {
	toParsed, err := ptypes.Timestamp(r.To)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	fromParsed, err := ptypes.Timestamp(r.From)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	sum, err := s.usecase.GetTotalUsageForCustomerInTimePeriod(fromParsed, toParsed, r.GetCustomerId())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	response := &timeseries_grpc.SumResponse{
		From:       r.From,
		To:         r.To,
		CustomerId: r.GetCustomerId(),
		Sum:        sum,
	}

	return response, nil
}
func (s *server) GetTotalUsageForMeterInTimePeriod(ctx context.Context, r *timeseries_grpc.GetRequest) (*timeseries_grpc.SumResponse, error) {
	toParsed, err := ptypes.Timestamp(r.To)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	fromParsed, err := ptypes.Timestamp(r.From)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	sum, err := s.usecase.GetTotalUsageForMeterInTimePeriod(fromParsed, toParsed, r.GetMeterId())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	response := &timeseries_grpc.SumResponse{
		From:    r.From,
		To:      r.To,
		MeterId: r.GetMeterId(),
		Sum:     sum,
	}

	return response, nil
}
