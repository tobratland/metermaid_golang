package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/tobra/metermaid/delivery/grpc/timeseries_grpc"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Client started...")

	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := timeseries_grpc.NewTimeseriesHandlerClient(conn)
	//store(client)
	//getdata(client)
	getSumCustomerId(client)

}

func getSumCustomerId(c timeseries_grpc.TimeseriesHandlerClient) {
	from, err := ptypes.TimestampProto(parseStringToTime("2018-08-05T00:00:00Z"))
	if err != nil {
		fmt.Println(err)
	}
	to, err := ptypes.TimestampProto(parseStringToTime("2018-08-20T23:00:00Z"))
	if err != nil {
		fmt.Println(err)
	}

	customerId := "tester123"

	request := &timeseries_grpc.GetRequest{
		From:       from,
		To:         to,
		CustomerId: customerId,
	}

	resp, err := c.GetTotalUsageForCustomerInTimePeriod(context.Background(), request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("sum for period ", resp.GetFrom(), " to period ", resp.GetTo(), "for customer: ", resp.GetCustomerId(), " is: ", resp.GetSum())
}

func getdata(c timeseries_grpc.TimeseriesHandlerClient) {
	from, err := ptypes.TimestampProto(parseStringToTime("2018-08-05T00:00:00Z"))
	if err != nil {
		fmt.Println(err)
	}
	to, err := ptypes.TimestampProto(parseStringToTime("2018-08-20T23:00:00Z"))
	if err != nil {
		fmt.Println(err)
	}
	request := &timeseries_grpc.GetRequest{
		From: from,
		To:   to,
	}

	resp, _ := c.GetAllFromTimeToTime(context.Background(), request)
	fmt.Println(len(resp.GetTimeseriesList()))
	for _, elem := range resp.GetTimeseriesList() {
		for _, value := range elem.Values {
			hour, err := ptypes.Timestamp(value.Hour)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(hour, ": ", value.Value)
		}
	}
}

func store(c timeseries_grpc.TimeseriesHandlerClient) {
	fromParsed, err := ptypes.TimestampProto(parseStringToTime("2018-08-09T00:00:00Z"))
	if err != nil {
		fmt.Println(err)
	}
	toParsed, err := ptypes.TimestampProto(parseStringToTime("2018-08-09T23:00:00Z"))
	if err != nil {
		fmt.Println(err)
	}
	request := &timeseries_grpc.Timeseries{
		MeterId:    "test123",
		CustomerId: "tester123",
		Resolution: "hour",
		From:       fromParsed,
		To:         toParsed,
	}

	req := addvalues(*request)

	resp, _ := c.Store(context.Background(), &req)
	fmt.Printf("Recieved responsse => [%s]", resp.CustomerId)
}

func parseStringToTime(s string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	res, err := time.Parse(layout, s)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func addvalues(t timeseries_grpc.Timeseries) timeseries_grpc.Timeseries {
	for i := 0; i < 24; i++ {
		str := fmt.Sprintf("2018-08-09T%02d:00:00Z", i)

		parsedTime, err := ptypes.TimestampProto(parseStringToTime(str))
		if err != nil {
			fmt.Println(err)
		}
		t.Values = append(t.Values, &timeseries_grpc.Value{
			Hour:  parsedTime,
			Value: 23.2,
		})
	}

	return t

}

/* Id         string             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
MeterId    string             `protobuf:"bytes,2,opt,name=meterId" json:"meterId,omitempty"`
CustomerId string             `protobuf:"bytes,3,opt,name=customerId" json:"customerId,omitempty"`
Resolution string             `protobuf:"bytes,4,opt,name=resolution" json:"resolution,omitempty"`
From       string             `protobuf:"bytes,5,opt,name=from" json:"from,omitempty"`
To         string             `protobuf:"bytes,6,opt,name=to" json:"to,omitempty"`
Values     map[string]float64 */
