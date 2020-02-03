package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func NewApiClient(url string, client *http.Client) *ApiClient {
	return &ApiClient{
		url:    url,
		client: client,
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
      "config_node_ip_address":"%s"
   }
}
`
	body = fmt.Sprintf(body, node.Name, node.IP)
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
	Name     string
	Hostname string
	IP       string
}
