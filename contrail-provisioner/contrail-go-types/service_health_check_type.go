//
// Automatically generated. DO NOT EDIT.
//

package types

type ServiceHealthCheckType struct {
	Enabled         bool   `json:"enabled,omitempty"`
	HealthCheckType string `json:"health_check_type,omitempty"`
	MonitorType     string `json:"monitor_type,omitempty"`
	Delay           int    `json:"delay,omitempty"`
	Delayusecs      int    `json:"delayUsecs,omitempty"`
	Timeout         int    `json:"timeout,omitempty"`
	Timeoutusecs    int    `json:"timeoutUsecs,omitempty"`
	MaxRetries      int    `json:"max_retries,omitempty"`
	HttpMethod      string `json:"http_method,omitempty"`
	UrlPath         string `json:"url_path,omitempty"`
	ExpectedCodes   string `json:"expected_codes,omitempty"`
}
