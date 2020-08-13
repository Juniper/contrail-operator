package keystone

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"k8s.io/client-go/rest"
)

func NewClient(config *rest.Config, k *contrail.Keystone) (*Client, error) {
	if k.Status.External {
		return &Client{client: NewDirectClient(k.Spec.ServiceConfiguration.AuthProtocol, k.Status.ClusterIP, k.Spec.ServiceConfiguration.ListenPort)}, nil
	}
	proxy, err := kubeproxy.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubeproxy: %v", err)
	}
	if k.Spec.ServiceConfiguration.AuthProtocol == "https" {
		return &Client{client: proxy.NewSecureClientForService(k.Namespace, k.Name+"-service", k.Status.Port)}, nil
	}
	return &Client{client: proxy.NewClientForService(k.Namespace, k.Name+"-service", k.Status.Port)}, nil
}

type KeystoneClient interface {
	NewRequest(method, path string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	client KeystoneClient
}

type DirectClient struct {
	url    string
	client http.Client
}

func NewDirectClient(protocol string, address string, port int) *DirectClient {
	url := fmt.Sprintf("%s://%s:%d", protocol, address, port)
	// Temporary don't validate certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &DirectClient{
		url:    url,
		client: http.Client{Transport: tr},
	}
}

func (c *DirectClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.url+path, body)
}

func (c *DirectClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) PostAuthTokens(username, password, project string) (AuthTokens, error) {
	return c.PostAuthTokensWithHeaders(username, password, project, http.Header{})
}

func (c *Client) PostAuthTokensWithHeaders(username, password, project string, headers http.Header) (AuthTokens, error) {
	kar := &keystoneAuthRequest{}
	kar.Auth.Identity.Methods = []string{"password"}
	kar.Auth.Identity.Password.User.Name = username
	kar.Auth.Identity.Password.User.Domain.ID = "default"
	kar.Auth.Identity.Password.User.Password = password
	kar.Auth.Scope.Project.Domain.ID = "default"
	kar.Auth.Scope.Project.Name = project
	karBody, err := json.Marshal(kar)
	if err != nil {
		return AuthTokens{}, err
	}
	request, err := c.client.NewRequest(http.MethodPost, "/v3/auth/tokens", bytes.NewReader(karBody))
	if err != nil {
		return AuthTokens{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	for name, values := range headers {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	response, err := c.client.Do(request)
	if err != nil {
		return AuthTokens{}, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return AuthTokens{}, newUnauthorized()
	}

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return AuthTokens{}, fmt.Errorf("invalid status code returned: %d", response.StatusCode)
	}

	authTokens := AuthTokens{}
	bytesRead, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return AuthTokens{}, err
	}
	if err := json.Unmarshal(bytesRead, &authTokens); err != nil {
		return AuthTokens{}, err
	}
	authTokens.XAuthTokenHeader = response.Header.Get("X-Subject-Token")
	return authTokens, nil
}

type keystoneAuthRequest struct {
	Auth struct {
		Identity struct {
			Methods  []string `json:"methods"`
			Password struct {
				User struct {
					Name   string `json:"name"`
					Domain struct {
						ID string `json:"id"`
					} `json:"domain"`
					Password string `json:"password"`
				} `json:"user"`
			} `json:"password"`
		} `json:"identity"`
		Scope struct {
			Project struct {
				Name   string `json:"name"`
				Domain struct {
					ID string `json:"id"`
				} `json:"domain"`
			} `json:"project"`
		} `json:"scope"`
	} `json:"auth"`
}

type AuthTokens struct {
	Token struct {
		Catalog []struct {
			Name      string
			Type      string
			Endpoints []struct {
				URL       string
				Interface string
			}
		}
	}
	XAuthTokenHeader string
}

func (t AuthTokens) EndpointURL(serviceName string, endpointInterface string) string {
	for _, service := range t.Token.Catalog {
		if service.Name == serviceName {
			for _, endpoint := range service.Endpoints {
				if endpoint.Interface == endpointInterface {
					return endpoint.URL
				}
			}
		}
	}
	return ""
}
