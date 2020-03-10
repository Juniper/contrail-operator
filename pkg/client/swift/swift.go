package swift

import (
	"bytes"
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

func (c *Client) PutContainer(name string) error {
	return c.PutContainerWithHeaders(name, http.Header{})
}

func (c *Client) PutReadAllContainer(name string) error {
	headers := http.Header{}
	headers.Add("X-Container-Read", ".r:*")
	return c.PutContainerWithHeaders(name, headers)
}

func (c *Client) PutContainerWithHeaders(name string, headers http.Header) error {
	request, err := c.proxy.NewRequest(http.MethodPut, c.path+"/"+name, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Auth-Token", c.token)
	for name, values := range headers {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	response, err := c.proxy.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
		return fmt.Errorf("invalid status code returned: %d, response: %s", response.StatusCode, c.response(response))
	}
	return nil
}

func (c *Client) GetContainer(name string) error {
	request, err := c.proxy.NewRequest(http.MethodGet, c.path+"/"+name, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Auth-Token", c.token)
	response, err := c.proxy.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
		return fmt.Errorf("invalid status code returned: %d, response: %s", response.StatusCode, c.response(response))
	}
	return nil
}

func (c *Client) response(response *http.Response) string {
	bodyAsString, _ := ioutil.ReadAll(response.Body)
	return string(bodyAsString)
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
		return fmt.Errorf("invalid status code returned: %d, response: %s", response.StatusCode, c.response(response))
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
		return nil, fmt.Errorf("invalid status code returned: %d, response: %s", response.StatusCode, c.response(response))
	}
	return ioutil.ReadAll(response.Body)
}

func (c *Client) GetFileWithoutAuth(container string, fileName string) ([]byte, error) {
	request, err := c.proxy.NewRequest(http.MethodGet, c.path+"/"+container+"/"+fileName, nil)
	if err != nil {
		return nil, err
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
