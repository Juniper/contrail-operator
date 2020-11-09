package config

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
)

type ConfigAPIResponse interface {
	IsValidConfigApiResponse() bool
}

type ConfigNodeResponse struct {
	Nodes []Node `json:"config-nodes"`
}

type VirtualRouterResponse struct {
	Nodes []Node `json:"virtual-routers"`
}

type BgpRouterResponse struct {
	Nodes []Node `json:"bgp-routers"`
}

type DatabaseNodeResponse struct {
	Nodes []Node `json:"database-nodes"`
}

type AnalyticsNodeResponse struct {
	Nodes []Node `json:"analytics-nodes"`
}

type Node struct {
	Href   string   `json:"href"`
	FqName []string `json:"fq_name"`
	UUID   string   `json:"uuid"`
}

func (c ConfigNodeResponse) IsValidConfigApiResponse() bool {
	if len(c.Nodes) > 0 {
		return true
	}
	return false
}

func NewClient(client *kubeproxy.Client, token string) (*Client, error) {
	return &Client{
		proxy: client,
		token: token,
	}, nil
}

type Client struct {
	proxy *kubeproxy.Client
	token string
}

func (c *Client) response(response *http.Response) string {
	bodyAsString, _ := ioutil.ReadAll(response.Body)
	return string(bodyAsString)
}

func (c *Client) GetResource(resourceName string) ([]byte, error) {
	request, err := c.proxy.NewRequest(http.MethodGet, resourceName, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		request.Header.Set("X-Auth-Token", c.token)
	}
	response, err := c.proxy.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code returned: %d, response: %s", response.StatusCode, c.response(response))
	}
	return ioutil.ReadAll(response.Body)
}
