package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tobra/metermaid/delivery/grpc/timeseries_grpc"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Client started...")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("127.0.0.1:8080", opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	client := timeseries_grpc.NewTimeseriesHandlerClient(cc)
	request := &timeseries_grpc.Timeseries{
		MeterId:    "test123",
		CustomerId: "tester123",
		Resolution: "hour",
		From:       "2018-08-09T00:00:00Z",
		To:         "2018-08-09T23:00:00Z",
		Values:     createMap(),
	}
	fmt.Println(request.Values["2018-08-09T00:00:00Z"])
	resp, _ := client.Store(context.Background(), request)
	fmt.Printf("Recieved responsse => [%s]", resp.Id)
}

func createMap() map[string]float64 {
	m := map[string]float64{
		"2018-08-09T00:00:00Z": 12.3,
		"2018-08-09T01:00:00Z": 12.3,
		"2018-08-09T02:00:00Z": 12.3,
		"2018-08-09T03:00:00Z": 12.3,
		"2018-08-09T04:00:00Z": 12.3,
		"2018-08-09T05:00:00Z": 12.3,
		"2018-08-09T06:00:00Z": 12.3,
		"2018-08-09T07:00:00Z": 12.3,
		"2018-08-09T08:00:00Z": 12.3,
		"2018-08-09T09:00:00Z": 12.3,
		"2018-08-09T10:00:00Z": 12.3,
		"2018-08-09T11:00:00Z": 12.3,
		"2018-08-09T12:00:00Z": 12.3,
		"2018-08-09T13:00:00Z": 12.3,
		"2018-08-09T14:00:00Z": 12.3,
		"2018-08-09T15:00:00Z": 12.3,
		"2018-08-09T16:00:00Z": 12.3,
		"2018-08-09T17:00:00Z": 12.3,
		"2018-08-09T18:00:00Z": 12.3,
		"2018-08-09T19:00:00Z": 12.3,
		"2018-08-09T20:00:00Z": 12.3,
		"2018-08-09T21:00:00Z": 12.3,
		"2018-08-09T22:00:00Z": 12.3,
		"2018-08-09T23:00:00Z": 12.3,
	}
	return m

}

/* Id         string             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
MeterId    string             `protobuf:"bytes,2,opt,name=meterId" json:"meterId,omitempty"`
CustomerId string             `protobuf:"bytes,3,opt,name=customerId" json:"customerId,omitempty"`
Resolution string             `protobuf:"bytes,4,opt,name=resolution" json:"resolution,omitempty"`
From       string             `protobuf:"bytes,5,opt,name=from" json:"from,omitempty"`
To         string             `protobuf:"bytes,6,opt,name=to" json:"to,omitempty"`
Values     map[string]float64 */
