package ml

import (
	"context"
	"encoding/json"
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

func (c OutlierCommand) Execute(_ context.Context, from, to time.Time, execute func(path string, payload []byte) ([]byte, error)) (*backend.QueryDataResponse, error) {
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
	return resp.Data, nil
}

func unmarshalOutlierCommand(data map[string]interface{}, interval time.Duration, appURL string) (*OutlierCommand, error) {
	uid, err := readValue[string](data, "datasource_uid")
	if err != nil {
		return nil, err
	}

	// TODO validate the data? What if ML API changes?
	/* data is expected to be like
			{
				"datasource_uid": "prometheus",
	            "datasource_type": "prometheus",
	            "query_params": {
	                "expr": "go_goroutines",
	                "range": true,
	                "refId": "A"
	            },
				"algorithm": {
	                "name": "dbscan",
	                "config": {
	                    "epsilon": 0.2354
	                },
	                "sensitivity": 0.9
	            },
	            "response_type": "binary"
			}
	*/

	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &OutlierCommand{
		query:         d,
		datasourceUID: uid,
		interval:      interval,
		appURL:        appURL,
	}, nil
}
