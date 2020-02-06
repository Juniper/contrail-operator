package swift

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Juniper/contrail-operator/test/kubeproxy"
)

func NewClient(client *kubeproxy.Client, token, endpointURL string) *Client {
	return &Client{
		proxy: client,
		token: token,
		path:  strings.TrimPrefix(endpointURL, "http://localhost:5080"),
	}
}

type Client struct {
	proxy *kubeproxy.Client
	token string
	path  string
}

func (c *Client) PutContainer(name string) error {
	request, err := c.proxy.NewRequest(http.MethodPut, c.path+"/"+name, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Auth-Token", c.token)
	response, err := c.proxy.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("invalid status code returned: %d", response.StatusCode)
	}
	return nil
}

func (c *Client) PutFile(container string, fileName string, content []byte) error {
	request, err := c.proxy.NewRequest(http.MethodPut, c.path+"/"+container+"/"+fileName, bytes.NewReader(content))
	if err != nil {
		return err
	}
	request.Header.Set("X-Auth-Token", c.token)
	response, err := c.proxy.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("invalid status code returned: %d", response.StatusCode)
	}
	return nil
}

func (c *Client) GetFile(container string, fileName string) ([]byte, error) {
	request, err := c.proxy.NewRequest(http.MethodGet, c.path+"/"+container+"/"+fileName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Auth-Token", c.token)
	response, err := c.proxy.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code returned: %d", response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}
