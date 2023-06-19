package datasources

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/components/simplejson"
)

func TestAllowedCookies(t *testing.T) {
	testCases := []struct {
		desc  string
		given map[string]interface{}
		want  AllowedCookies
		// remove this then FlagAllowedCookieRegexPattern is removed
		isFlagAllowedCookieRegexPatternEnabled bool
	}{
		{
			desc: "Usual json data without any pattern matching option provided",
			given: map[string]interface{}{
				"keepCookies": []string{"cookie2"},
			},
			want: AllowedCookies{
				MatchOption:  MO_EXACT_MATCH,
				MatchPattern: "",
				KeepCookies:  []string{"cookie2"},
			},
		},
		{
			desc: "Usual json data with pattern matching option provided",
			given: map[string]interface{}{
				"keepCookies":          []string{"cookie2"},
				"allowedCookieOption":  "exact_matching",
				"allowedCookiePattern": "",
			},
			want: AllowedCookies{
				MatchOption:  MO_EXACT_MATCH,
				MatchPattern: "",
				KeepCookies:  []string{"cookie2"},
			},
		},
		{
			desc: "Usual json data with unknown pattern matching option provided",
			given: map[string]interface{}{
				"keepCookies":          []string{"cookie2"},
				"allowedCookieOption":  "unknown_option",
				"allowedCookiePattern": "",
			},
			want: AllowedCookies{
				MatchOption:  MO_EXACT_MATCH,
				MatchPattern: "",
				KeepCookies:  []string{"cookie2"},
			},
		},
		{
			desc: "Usual json data with empty pattern matching option provided",
			given: map[string]interface{}{
				"keepCookies":          []string{"cookie2"},
				"allowedCookieOption":  "",
				"allowedCookiePattern": "",
			},
			want: AllowedCookies{
				MatchOption:  MO_EXACT_MATCH,
				MatchPattern: "",
				KeepCookies:  []string{"cookie2"},
			},
		},
		{
			desc: "Usual json data with regex pattern matching option provided",
			given: map[string]interface{}{
				"keepCookies":          []string{"cookie2"},
				"allowedCookieOption":  "regex_match",
				"allowedCookiePattern": ".*",
			},
			want: AllowedCookies{
				MatchOption:  MO_REGEX_MATCH,
				MatchPattern: ".*",
				KeepCookies:  []string{"cookie2"},
			},
			isFlagAllowedCookieRegexPatternEnabled: true,
		},
		{
			desc: "Json data with regex pattern matching option provided",
			given: map[string]interface{}{
				"keepCookies":          []string{"cookie2"},
				"allowedCookieOption":  "regex_match",
				"allowedCookiePattern": `\w+`,
			},
			want: AllowedCookies{
				MatchOption:  MO_REGEX_MATCH,
				MatchPattern: `\w+`,
				KeepCookies:  []string{"cookie2"},
			},
			isFlagAllowedCookieRegexPatternEnabled: true,
		},
		{
			desc: "Json data with regex pattern matching option provided and empty keepCookies",
			given: map[string]interface{}{
				"keepCookies":          []string{},
				"allowedCookieOption":  "regex_match",
				"allowedCookiePattern": `^special_.*`,
			},
			want: AllowedCookies{
				MatchOption:  MO_REGEX_MATCH,
				MatchPattern: `^special_.*`,
				KeepCookies:  []string{},
			},
			isFlagAllowedCookieRegexPatternEnabled: true,
		},
		{
			desc: "Json data with regex pattern matching option provided and no keepCookies",
			given: map[string]interface{}{
				"allowedCookieOption":  "regex_match",
				"allowedCookiePattern": `^special_.*`,
			},
			want: AllowedCookies{
				MatchOption:  MO_REGEX_MATCH,
				MatchPattern: `^special_.*`,
				KeepCookies:  nil,
			},
			isFlagAllowedCookieRegexPatternEnabled: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			jsonDataBytes, err := json.Marshal(&test.given)
			require.NoError(t, err)
			jsonData, err := simplejson.NewJson(jsonDataBytes)
			require.NoError(t, err)

			ds := DataSource{
				ID:       1235,
				JsonData: jsonData,
				UID:      "test",
			}

			actual := ds.AllowedCookies(test.isFlagAllowedCookieRegexPatternEnabled)
			assert.Equal(t, test.want.MatchOption, actual.MatchOption)
			assert.Equal(t, test.want.MatchPattern, actual.MatchPattern)
			assert.EqualValues(t, test.want.KeepCookies, actual.KeepCookies)
		})
	}
}
