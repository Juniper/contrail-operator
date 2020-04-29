//
// Automatically generated. DO NOT EDIT.
//

package types

type LoadbalancerHealthmonitorType struct {
	AdminState bool `json:"admin_state,omitempty"`
	MonitorType string `json:"monitor_type,omitempty"`
	Delay int `json:"delay,omitempty"`
	Timeout int `json:"timeout,omitempty"`
	MaxRetries int `json:"max_retries,omitempty"`
	HttpMethod string `json:"http_method,omitempty"`
	UrlPath string `json:"url_path,omitempty"`
	ExpectedCodes string `json:"expected_codes,omitempty"`
}
