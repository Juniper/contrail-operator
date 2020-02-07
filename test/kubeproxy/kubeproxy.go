package kubeproxy

import (
	"fmt"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func New(config *rest.Config) (*HTTPProxy, error) {
	transport, err := rest.TransportFor(config)
	if err != nil {
		return nil, err
	}
	client := http.Client{Transport: transport}
	url, s, err := rest.DefaultServerURL(config.Host, config.APIPath, schema.GroupVersion{}, true)
	if err != nil {
		return nil, err
	}
	return &HTTPProxy{
		client:    client,
		serverURL: url.String() + s,
	}, nil
}

type HTTPProxy struct {
	client    http.Client
	serverURL string
}

func (p *HTTPProxy) NewClient(namespace, pod string, port int) *Client {
	return p.NewClientWithPath(namespace, pod, port, "")
}

func (p *HTTPProxy) NewClientWithPath(namespace, pod string, port int, path string) *Client {
	url := fmt.Sprintf("%sapi/v1/namespaces/%s/pods/%s:%d/proxy%s", p.serverURL, namespace, pod, port, path)
	return &Client{
		url:    url,
		client: p.client,
	}
}

// Client is an HTTP client using Kubernetes API server as a proxy. With this client you can execute HTTP methods
// on any k8s pod from outside the cluster. No need to manually execute `kubectl port-forward`, `kubectl proxy`
// or something similar.
type Client struct {
	url    string
	client http.Client
}

func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.url+path, body)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	fmt.Println("DO URL:", req.URL)
	return c.client.Do(req)
}
