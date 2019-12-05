package config

import (
	"bytes"
	"fmt"
	"net/http"
	"errors"
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
   "config-node":{ 
      "parent_type":"global-system-config",
      "fq_name":[ 
         "default-global-system-config",
         "%s"
      ],
      "uuid":"%s",
      "config_node_ip_address":"%s"
   }
}
`
	body = fmt.Sprintf(body, node.Hostname, node.UUID, node.IP)
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
	UUID     string
	Hostname string
	IP       string
}
