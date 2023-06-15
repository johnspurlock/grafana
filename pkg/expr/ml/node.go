package ml

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)


const (
	Outlier CommandType = "outlier"

	// format of the time used by outlier API
	timeFormat = "2006-01-02T15:04:05.999999999"

	defaultInterval = 1000 * time.Millisecond
)

type Command interface {
	DatasourceUID() string
	Execute(ctx context.Context, from, to time.Time, executor func(path string, payload []byte) ([]byte, error)) (*backend.QueryDataResponse, error)
}

type CommandType string


func UnmarshalCommand(query map[string]interface{}, appURL string) (Command, error) {
	mlType, err := readValue[string](query, "type")
	if err != nil {
		return nil, err
	}

	intervalMs, err := readOptionalValue[float64](query, "intervalMs")
	if err != nil {
		return nil, err
	}
	interval := defaultInterval
	if intervalMs != nil {
		interval = time.Duration(*intervalMs) * time.Millisecond
	}

	d, err := readValue[map[string]interface{}](query, "data")
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(mlType) {
	case string(Outlier):
		return unmarshalOutlierCommand(d, interval, appURL)
	default:
		return nil, fmt.Errorf("unsupported command type '%v'. Supported only '%s'", mlType, Outlier)
	}
}
