package datasources

import (
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/services/quota"
	"github.com/grafana/grafana/pkg/services/user"
)

const (
	DS_GRAPHITE       = "graphite"
	DS_INFLUXDB       = "influxdb"
	DS_INFLUXDB_08    = "influxdb_08"
	DS_ES             = "elasticsearch"
	DS_PROMETHEUS     = "prometheus"
	DS_ALERTMANAGER   = "alertmanager"
	DS_JAEGER         = "jaeger"
	DS_LOKI           = "loki"
	DS_OPENTSDB       = "opentsdb"
	DS_TEMPO          = "tempo"
	DS_ZIPKIN         = "zipkin"
	DS_MYSQL          = "mysql"
	DS_POSTGRES       = "postgres"
	DS_MSSQL          = "mssql"
	DS_ACCESS_DIRECT  = "direct"
	DS_ACCESS_PROXY   = "proxy"
	DS_ES_OPEN_DISTRO = "grafana-es-open-distro-datasource"
	DS_ES_OPENSEARCH  = "grafana-opensearch-datasource"
	DS_AZURE_MONITOR  = "grafana-azure-monitor-datasource"
	// CustomHeaderName is the prefix that is used to store the name of a custom header.
	CustomHeaderName = "httpHeaderName"
	// CustomHeaderValue is the prefix that is used to store the value of a custom header.
	CustomHeaderValue = "httpHeaderValue"
	// MO_EXACT_MATCH is for the allowed cookies. Cookie name must match exactly.
	MO_EXACT_MATCH = "exact_match"
	// MO_REGEX_MATCH is for the allowed cookies. Cookie name must match with the provided regex.
	MO_REGEX_MATCH = "regex_match"
)

type DsAccess string

type DataSource struct {
	ID      int64 `json:"id,omitempty" xorm:"pk autoincr 'id'"`
	OrgID   int64 `json:"orgId,omitempty" xorm:"org_id"`
	Version int   `json:"version,omitempty"`

	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Access DsAccess `json:"access"`
	URL    string   `json:"url" xorm:"url"`
	// swagger:ignore
	Password      string `json:"-"`
	User          string `json:"user"`
	Database      string `json:"database"`
	BasicAuth     bool   `json:"basicAuth"`
	BasicAuthUser string `json:"basicAuthUser"`
	// swagger:ignore
	BasicAuthPassword string            `json:"-"`
	WithCredentials   bool              `json:"withCredentials"`
	IsDefault         bool              `json:"isDefault"`
	JsonData          *simplejson.Json  `json:"jsonData"`
	SecureJsonData    map[string][]byte `json:"secureJsonData"`
	ReadOnly          bool              `json:"readOnly"`
	UID               string            `json:"uid" xorm:"uid"`

	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

type AllowedCookies struct {
	MatchOption  string
	MatchPattern string
	KeepCookies  []string
}

// AllowedCookies parses the jsondata.KeepCookies and returns a list of
// allowed cookies, otherwise an empty list.
func (ds DataSource) AllowedCookies(isFlagAllowedCookieRegexPattern bool) AllowedCookies {
	// Default matching option is exact matching
	ac := AllowedCookies{
		MatchOption:  MO_EXACT_MATCH,
		MatchPattern: "",
		KeepCookies:  []string{},
	}
	if ds.JsonData != nil {
		aco := ds.JsonData.Get("allowedCookieOption")
		acp := ds.JsonData.Get("allowedCookiePattern")
		kc := ds.JsonData.Get("keepCookies")

		if kc != nil {
			ac.KeepCookies = kc.MustStringArray()
		}

		// Remove this check when FlagAllowedCookieRegexPattern is removed
		if isFlagAllowedCookieRegexPattern {
			if acp != nil {
				ac.MatchPattern = acp.MustString()
			}

			if aco != nil {
				acoStr := aco.MustString()
				if acoStr == MO_EXACT_MATCH || acoStr == MO_REGEX_MATCH {
					ac.MatchOption = acoStr
				}
			}
		} else {
			ac.MatchPattern = ""
			ac.MatchOption = MO_EXACT_MATCH
		}
	}

	return ac
}

// Specific error type for grpc secrets management so that we can show more detailed plugin errors to users
type ErrDatasourceSecretsPluginUserFriendly struct {
	Err string
}

func (e ErrDatasourceSecretsPluginUserFriendly) Error() string {
	return e.Err
}

// ----------------------
// COMMANDS

// Also acts as api DTO
type AddDataSourceCommand struct {
	Name            string            `json:"name" binding:"Required"`
	Type            string            `json:"type" binding:"Required"`
	Access          DsAccess          `json:"access" binding:"Required"`
	URL             string            `json:"url"`
	Database        string            `json:"database"`
	User            string            `json:"user"`
	BasicAuth       bool              `json:"basicAuth"`
	BasicAuthUser   string            `json:"basicAuthUser"`
	WithCredentials bool              `json:"withCredentials"`
	IsDefault       bool              `json:"isDefault"`
	JsonData        *simplejson.Json  `json:"jsonData"`
	SecureJsonData  map[string]string `json:"secureJsonData"`
	UID             string            `json:"uid"`

	OrgID                   int64             `json:"-"`
	UserID                  int64             `json:"-"`
	ReadOnly                bool              `json:"-"`
	EncryptedSecureJsonData map[string][]byte `json:"-"`
	UpdateSecretFn          UpdateSecretFn    `json:"-"`
}

// Also acts as api DTO
type UpdateDataSourceCommand struct {
	Name            string            `json:"name" binding:"Required"`
	Type            string            `json:"type" binding:"Required"`
	Access          DsAccess          `json:"access" binding:"Required"`
	URL             string            `json:"url"`
	User            string            `json:"user"`
	Database        string            `json:"database"`
	BasicAuth       bool              `json:"basicAuth"`
	BasicAuthUser   string            `json:"basicAuthUser"`
	WithCredentials bool              `json:"withCredentials"`
	IsDefault       bool              `json:"isDefault"`
	JsonData        *simplejson.Json  `json:"jsonData"`
	SecureJsonData  map[string]string `json:"secureJsonData"`
	Version         int               `json:"version"`
	UID             string            `json:"uid"`

	OrgID                   int64             `json:"-"`
	ID                      int64             `json:"-"`
	ReadOnly                bool              `json:"-"`
	EncryptedSecureJsonData map[string][]byte `json:"-"`
	UpdateSecretFn          UpdateSecretFn    `json:"-"`
}

// DeleteDataSourceCommand will delete a DataSource based on OrgID as well as the UID (preferred), ID, or Name.
// At least one of the UID, ID, or Name properties must be set in addition to OrgID.
type DeleteDataSourceCommand struct {
	ID   int64
	UID  string
	Name string

	OrgID int64

	DeletedDatasourcesCount int64

	UpdateSecretFn UpdateSecretFn
}

// Function for updating secrets along with datasources, to ensure atomicity
type UpdateSecretFn func() error

// ---------------------
// QUERIES

type GetDataSourcesQuery struct {
	OrgID           int64
	DataSourceLimit int
	User            *user.SignedInUser
}

type GetAllDataSourcesQuery struct{}

type GetDataSourcesByTypeQuery struct {
	OrgID int64 // optional: filter by org_id
	Type  string
}

type GetDefaultDataSourceQuery struct {
	OrgID int64
	User  *user.SignedInUser
}

// GetDataSourceQuery will get a DataSource based on OrgID as well as the UID (preferred), ID, or Name.
// At least one of the UID, ID, or Name properties must be set in addition to OrgID.
type GetDataSourceQuery struct {
	ID   int64
	UID  string
	Name string

	OrgID int64
}

type DatasourcesPermissionFilterQuery struct {
	User        *user.SignedInUser
	Datasources []*DataSource
}

const (
	QuotaTargetSrv quota.TargetSrv = "data_source"
	QuotaTarget    quota.Target    = "data_source"
)
