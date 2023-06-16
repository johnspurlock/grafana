package ml

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/util/cmputil"
)

func TestOutlierExec(t *testing.T) {
	outlier := OutlierCommand{
		query: jsoniter.RawMessage(`
		{
			"datasource_uid": "a4ce599c-4c93-44b9-be5b-76385b8c01be",
			"datasource_type": "prometheus",
			"query_params": {
				"expr": "go_goroutines{}",
				"range": true,
				"refId": "A"
			},
			"response_type": "binary",
			"algorithm": {
				"name": "dbscan",
				"config": {
					"epsilon": 7.667
				},
				"sensitivity": 0.83
			}
		}	
		`),
		datasourceUID: "a4ce599c-4c93-44b9-be5b-76385b8c01be",
		appURL:        "https://grafana.com",
		interval:      1000 * time.Second,
	}

	successResponse := `{"status":"success","data":{"results":{"A":{"status":200,"frames":[{"schema":{"name":"go_goroutines{instance=\"localhost:9090\", job=\"prometheus\"}","fields":[{"name":"Time","type":"time","typeInfo":{"frame":"time.Time"}},{"name":"Value","type":"number","typeInfo":{"frame":"float64"},"labels":{"__name__":"go_goroutines","instance":"localhost:9090","job":"prometheus"},"config":{"displayNameFromDS":"go_goroutines{instance=\"localhost:9090\", job=\"prometheus\"}"}}]},"data":{"values":[[1686945300000,1686945360000,1686945420000,1686945480000,1686945540000],[0,0,0,0,0]]}},{"schema":{"name":"go_goroutines{instance=\"mlapi:4030\", job=\"mlapi\"}","fields":[{"name":"Time","type":"time","typeInfo":{"frame":"time.Time"}},{"name":"Value","type":"number","typeInfo":{"frame":"float64"},"labels":{"__name__":"go_goroutines","instance":"mlapi:4030","job":"mlapi"},"config":{"displayNameFromDS":"go_goroutines{instance=\"mlapi:4030\", job=\"mlapi\"}"}}]},"data":{"values":[[1686945300000,1686945360000,1686945420000,1686945480000,1686945540000],[0,0,0,0,0]]}},{"schema":{"name":"go_goroutines{instance=\"scheduler:4031\", job=\"scheduler\"}","fields":[{"name":"Time","type":"time","typeInfo":{"frame":"time.Time"}},{"name":"Value","type":"number","typeInfo":{"frame":"float64"},"labels":{"__name__":"go_goroutines","instance":"scheduler:4031","job":"scheduler"},"config":{"displayNameFromDS":"go_goroutines{instance=\"scheduler:4031\", job=\"scheduler\"}"}}]},"data":{"values":[[1686945300000,1686945360000,1686945420000,1686945480000,1686945540000],[0,0,0,0,0]]}},{"schema":{"name":"go_goroutines{instance=\"sift:4032\", job=\"sift\"}","fields":[{"name":"Time","type":"time","typeInfo":{"frame":"time.Time"}},{"name":"Value","type":"number","typeInfo":{"frame":"float64"},"labels":{"__name__":"go_goroutines","instance":"sift:4032","job":"sift"},"config":{"displayNameFromDS":"go_goroutines{instance=\"sift:4032\", job=\"sift\"}"}}]},"data":{"values":[[1686945300000,1686945360000,1686945420000,1686945480000,1686945540000],[0,0,0,0,0]]}}]}}}}`

	var templateMap map[string]interface{}
	require.NoError(t, json.Unmarshal(outlier.query, &templateMap))

	t.Run("should create correct payload and parse data frames", func(t *testing.T) {
		to := time.Now()
		from := to.Add(-10 * time.Hour)

		resp, err := outlier.Execute(from, to, func(path string, payload []byte) ([]byte, error) {
			require.Equal(t, "/proxy/api/v1/outlier", path)
			var payloadMap map[string]interface{}
			payloadAny := jsoniter.Get(payload, "data", "attributes").MustBeValid()
			payloadAny.ToVal(&payloadMap)

			var reporter cmputil.DiffReporter
			cmp.Diff(templateMap, payloadMap, cmp.Reporter(&reporter))

			require.ElementsMatch(t, []string{
				"[start_end_attributes]",
				"[grafana_url]",
			}, reporter.Diffs.Paths())

			require.EqualValues(t, from.Format(timeFormat), payloadAny.Get("start_end_attributes", "start").ToString())
			require.EqualValues(t, to.Format(timeFormat), payloadAny.Get("start_end_attributes", "end").ToString())
			require.EqualValues(t, outlier.interval.Milliseconds(), payloadAny.Get("start_end_attributes", "interval").ToInt64())

			require.EqualValues(t, outlier.appURL, payloadAny.Get("grafana_url").ToString())

			return []byte(successResponse), nil
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Responses, 1)
		require.NotEmpty(t, resp.Responses["A"].Frames[0].Rows())
	})
}
