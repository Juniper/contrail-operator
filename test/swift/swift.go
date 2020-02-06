package swift

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Juniper/contrail-operator/test/kubeproxy"
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
		all, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(all))
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
