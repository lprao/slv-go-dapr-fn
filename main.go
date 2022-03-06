package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/lprao/slv-go-lib/pkg/logger"
	"github.com/lprao/slv-go-lib/pkg/slvlib"
	slvpb "github.com/lprao/slv-proto"
)

const (
	SensorName = "Sensor#93e00902"
)

type SensorValue struct {
	MoistureLevel int `json: "moistureLevel"`
}

func main() {
	s := daprd.NewService(":8083")

	if err := s.AddServiceInvocationHandler("/", sensorHandler); err != nil {
		Logger.Fatalf("error adding invocation handler: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		Logger.Fatalf("error listenning: %v", err)
	}
}

func sensorHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	var sensorValue SensorValue
	var slv *slvlib.SlvInt

	if in == nil {
		err = errors.New("invocation parameter required")
		return
	}

	if err = json.Unmarshal(in.Data, &sensorValue); err != nil {
		return
	}

	c, err := client.NewClient()
	if err != nil {
		return
	}

	opt := map[string]string{
		"version": "2",
	}
	slvVarName, err := c.GetSecret(ctx, "dapr-secret-store", SensorName, opt)
	if err != nil {
		return
	}

	slv, err = slvlib.GetSlvIntByName(slvVarName[SensorName])
	if err != nil {
		slv, err = slvlib.NewSlvInt(slvVarName[SensorName], 0, slvpb.VarScope_PRIVATE, slvpb.VarPermissions_READWRITE)
		if err != nil {
			return
		}
	}

	_, err = slv.Set(sensorValue.MoistureLevel)
	if err != nil {
		return
	}

	out = &common.Content{
		Data:        []byte("Updates Sensor value"),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}

var Logger *logger.Log

func init() {
	Logger = logger.NewLogger()
}
