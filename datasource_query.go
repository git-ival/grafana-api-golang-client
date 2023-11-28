package gapi

import (
	jsonx "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type DataSourceQueryType string

// non-exhaustive
const (
	QueryTypeTimeSeries DataSourceQueryType = "timeSeriesQuery"
)

// Example: {From: "2023-11-11T18:39:16.966Z", To: "2023-11-11T19:39:16.966Z", Raw: {From: "now-1h", To: "now"}}
type QueryTimeRange struct {
	From string                      `json:"from,omitempty"` // RFC3339
	To   string                      `json:"to,omitempty"`   // RFC3339
	Raw  *DashboardRelativeTimeRange `json:"raw"`
}

type DataSourceQueryQuery struct {
	DataSource    DataSource     `json:"datasource"`
	Format        string         `json:"format,omitempty"`
	Expression    string         `json:"expr,omitempty"`
	Interval      string         `json:"interval,omitempty"`
	LegendFormat  string         `json:"legendFormat,omitempty"`
	RefID         string         `json:"refId,omitempty"`
	RequestID     string         `json:"requestId,omitempty`
	ScenarioID    string         `json:"scenarioId,omitempty"`
	QueryType     string         `json:"queryType,omitempty"`
	Exemplar      bool           `json:"exemplar,omitempty"`
	UTCOffsetSec  int64          `json:"utcOffsetSec,omitempty"`
	DatasourceID  int64          `json:"datasourceId,omitempty"`
	IntervalMs    int64          `json:"intervalMs,omitempty"`
	MaxDataPoints int64          `json:"maxDataPoints,omitempty"`
	StringInput   string         `json:"stringInput,omitempty"`
	Unknown       jsontext.Value `json:",unknown"`
}

type DataSourceQuery struct {
	Queries []DataSourceQueryQuery `json:"queries,omitempty"`
	Range   QueryTimeRange         `json:"range,omitempty"`
	From    string                 `json:"from,omitempty"` // Unix Milliseconds
	To      string                 `json:"to,omitempty"`   // Unix Milliseconds
}

type Field struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
	Labels map[string]string      `json:"labels,omitempty"`
}

type SchemaField struct {
	Field
	TypeInfo map[string]interface{} `json:"typeInfo"`
}

type Schema struct {
	Name  string `json:"name,omitempty"`
	RefID string `json:"refId"`
	Meta  *struct {
		Custom *struct {
			ResultType string `json:"resultType,omitempty"`
		} `json:"custom,omitempty"`
		ExecutedQueryString        string `json:"executedQueryString,omitempty"`
		PreferredVisualisationType string `json:"preferredVisualisationType,omitempty"`
		Type                       string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
	Fields  []SchemaField  `json:"fields"`
	Unknown jsontext.Value `json:",unknown"`
}

type QueryResult struct {
	Frames []struct {
		Schema Schema `json:"schema"`
		Data   struct {
			Values [][]interface{} `json:"values"`
		} `json:"data"`
	} `json:"frames"`
	Unknown jsontext.Value `json:",unknown"`
}

type QueryResults map[string]*QueryResult

type DataSourceQueryResponse struct {
	Results QueryResults `json:"results"`
}

func (c *Client) SchemaFieldToSnapshotField(schemaField SchemaField) SnapshotField {
	return SnapshotField{
		Field: Field{
			Config: schemaField.Config,
			Labels: schemaField.Labels,
			Name:   schemaField.Name,
			Type:   schemaField.Type,
		},
		Values: []interface{}{},
	}
}

func (c *Client) QueryDataSource(dsq DataSourceQuery) (*QueryResults, error) {
	path := "/api/ds/query"

	data, err := jsonx.Marshal(dsq, defaultJSONOptions()...)
	if err != nil {
		return nil, err
	}

	result := &DataSourceQueryResponse{}

	err = c.request("POST", path, nil, data, &result)
	if err != nil {
		return nil, err
	}

	return &result.Results, err
}
