package ml

import (
	"errors"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	jsoniter "github.com/json-iterator/go"
)

type OutlierCommand struct {
	query         jsoniter.RawMessage
	datasourceUID string
	appURL        string
	interval      time.Duration
}

func (c OutlierCommand) DatasourceUID() string {
	return c.datasourceUID
}

func (c OutlierCommand) Execute(from, to time.Time, execute func(path string, payload []byte) ([]byte, error)) (*backend.QueryDataResponse, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal(c.query, &dataMap)
	if err != nil {
		return nil, err
	}

	dataMap["start_end_attributes"] = map[string]interface{}{
		"start":    from.Format(timeFormat),
		"end":      to.Format(timeFormat),
		"interval": c.interval.Milliseconds(),
	}
	dataMap["grafana_url"] = c.appURL

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": dataMap,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	responseBody, err := execute("/proxy/api/v1/outlier", body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Status string                     `json:"status"`
		Data   *backend.QueryDataResponse `json:"data,omitempty"`
		Error  string                     `json:"error,omitempty"`
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("cannot umarshall response from plugin API: %w", err)
	}
	if len(resp.Error) > 0 {
		return nil, errors.New(resp.Error)
	}
	return resp.Data, nil
}

func unmarshalOutlierCommand(q jsoniter.Any, appURL string) (*OutlierCommand, error) {
	interval := defaultInterval
	switch intervalNode := q.Get("intervalMs"); intervalNode.ValueType() {
	case jsoniter.NilValue:
	case jsoniter.InvalidValue:
	case jsoniter.NumberValue:
		interval = time.Duration(intervalNode.ToInt64()) * time.Millisecond
	default:
		return nil, fmt.Errorf("field `intervalMs` is expected to be a number")
	}

	cfgNode := q.Get("config")
	if cfgNode.ValueType() != jsoniter.ObjectValue {
		return nil, fmt.Errorf("field `config` is required and should be object")
	}
	ds := cfgNode.Get("datasource_uid").ToString()
	if len(ds) == 0 {
		return nil, fmt.Errorf("field `config.datasource_uid` is required and should be string")
	}

	d, err := json.Marshal(cfgNode.GetInterface())
	if err != nil {
		return nil, err
	}

	return &OutlierCommand{
		query:         d,
		datasourceUID: ds,
		interval:      interval,
		appURL:        appURL,
	}, nil
}
