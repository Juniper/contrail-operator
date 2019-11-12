package config

import (
	"bytes"
	"fmt"
	"net/http"

	"atom/atom/logging-service/errors"
)

func NewApiClient(url string) *ApiClient {
	return &ApiClient{
		url:    url,
		client: &http.Client{},
	}
}

// ApiClient is a client for communicating with the Config API server
type ApiClient struct {
	url    string
	client *http.Client
}

func (c *ApiClient) EnsureConfigNodeExists(node ConfigNode) error {
	body := `
{
    "config-node": {
        "fq_name": [
            "default-global-system-config",
            "%s"
        ],
        "uuid": "%s"
    }
}
`
	body = fmt.Sprintf(body, node.Host, node.UUID)
	response, err := c.client.Post(c.url+"/config-nodes", "application/json", bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	if configNodeAlreadyExists := response.StatusCode == 409; configNodeAlreadyExists {
		return nil
	}
	if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("error status code received: %d", response.StatusCode))
	}
	return nil
}

type ConfigNode struct {
	UUID string
	Host string
}
