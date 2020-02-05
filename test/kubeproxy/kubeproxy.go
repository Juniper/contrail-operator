package kubeproxy

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func New(t *testing.T, config *rest.Config) *HTTPProxy {
	transport, err := rest.TransportFor(config)
	require.NoError(t, err)
	client := http.Client{Transport: transport}
	url, s, err := rest.DefaultServerURL(config.Host, config.APIPath, schema.GroupVersion{}, true)
	require.NoError(t, err)

	return &HTTPProxy{
		t:         t,
		client:    client,
		serverURL: url.String() + s,
	}
}

type HTTPProxy struct {
	client    http.Client
	serverURL string
	t         *testing.T
}

func (p *HTTPProxy) ClientFor(namespace, pod string, port int) *Client {
	url := fmt.Sprintf("%sapi/v1/namespaces/%s/pods/%s:%d/proxy", p.serverURL, namespace, pod, port)
	return &Client{
		url:    url,
		client: p.client,
		t:      p.t,
	}
}

// Client is an HTTP client using Kubernetes API server as a proxy. With this client you can execute HTTP methods
// on any k8s pod from outside the cluster. No need to manually execute `kubectl port-forward`, `kubectl proxy`
// or something similar.
type Client struct {
	url    string
	client http.Client
	t      *testing.T
}

func (c *Client) NewRequest(method, path string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, c.url+path, body)
	require.NoError(c.t, err)
	return request
}

func (c *Client) Do(req *http.Request) *http.Response {
	response, err := c.client.Do(req)
	require.NoError(c.t, err)
	return response
}
