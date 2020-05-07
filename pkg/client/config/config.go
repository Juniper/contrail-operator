package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
)

func NewClient(client *kubeproxy.Client, token, endpointURL string) (*Client, error) {
	fullURL, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		proxy: client,
		token: token,
		path:  fullURL.Path,
	}, nil
}

type Client struct {
	proxy *kubeproxy.Client
	token string
	path  string
}

func (c *Client) response(response *http.Response) string {
	bodyAsString, _ := ioutil.ReadAll(response.Body)
	return string(bodyAsString)
}

func (c *Client) GetResource(resourceName string) ([]byte, error) {
	request, err := c.proxy.NewRequest(http.MethodGet, c.path+resourceName, nil)
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
