package swift

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Juniper/contrail-operator/test/kubeproxy"
)

func NewClient(t *testing.T, client *kubeproxy.Client, token, endpointURL string) *Client {
	return &Client{
		proxy: client,
		t:     t,
		token: token,
		path:  strings.TrimPrefix(endpointURL, "http://localhost:5080"),
	}
}

type Client struct {
	proxy *kubeproxy.Client
	t     *testing.T
	token string
	path  string
}

func (c *Client) PutContainer(name string) {
	request := c.proxy.NewRequest(http.MethodPut, c.path+"/"+name, nil)
	request.Header.Set("X-Auth-Token", c.token)
	response := c.proxy.Do(request)
	require.Equal(c.t, 201, response.StatusCode)
}

func (c *Client) PutFile(container string, fileName string, content []byte) {
	request := c.proxy.NewRequest(http.MethodPut, c.path+"/"+container+"/"+fileName, bytes.NewReader(content))
	request.Header.Set("X-Auth-Token", c.token)
	response := c.proxy.Do(request)
	require.Equal(c.t, 201, response.StatusCode)
}

func (c *Client) GetFile(container string, fileName string) []byte {
	request := c.proxy.NewRequest(http.MethodGet, c.path+"/"+container+"/"+fileName, nil)
	request.Header.Set("X-Auth-Token", c.token)
	response := c.proxy.Do(request)
	require.Equal(c.t, 200, response.StatusCode)
	content, err := ioutil.ReadAll(response.Body)
	require.NoError(c.t, err)
	return content
}
