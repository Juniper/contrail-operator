package keystone

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
)

// NewClient prepares keystone client based on properties of passed keystone CRD.
func NewClient(kubClient client.Client, scheme *runtime.Scheme, config *rest.Config, k *contrail.Keystone) (*Client, error) {
	connector, err := newConnector(kubClient, scheme, config, k)
	if err != nil {
		return nil, err
	}
	return &Client{
		Connector:    connector,
		KeystoneConf: &k.Spec.ServiceConfiguration,
	}, nil
}

func newConnector(kubClient client.Client, scheme *runtime.Scheme, config *rest.Config, k *contrail.Keystone) (keystoneClient, error) {
	if k.Spec.ServiceConfiguration.ExternalAddress != "" {
		caCertificate := certificates.NewCACertificate(kubClient, scheme, k, k.GetName())
		caBundle, _ := caCertificate.GetCaCert()
		return newExtKeystoneClient(k.Spec.ServiceConfiguration.AuthProtocol, k.Spec.ServiceConfiguration.ExternalAddress, k.Spec.ServiceConfiguration.ListenPort, caBundle), nil
	}
	proxy, err := kubeproxy.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubeproxy: %v", err)
	}
	if k.Spec.ServiceConfiguration.AuthProtocol == "https" {
		return proxy.NewSecureClientForService(k.Namespace, k.Name+"-service", k.Status.Port), nil
	}
	return proxy.NewClientForService(k.Namespace, k.Name+"-service", k.Status.Port), nil
}

type keystoneClient interface {
	NewRequest(method, path string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}

// A Client is an interface to the Keystone endpoint which allows retrieving
// tokens, endpoints etc.
type Client struct {
	// Connector specifies backend mechanism used to communicate with Keystone.
	// When keystone is deployed as part of the cluster this should be a kubeproxy client.
	// If keystone service resides outside of the cluster, then general http client
	// can be used which implements keystoneClient interface methods.
	Connector keystoneClient
	// Service configuration of the Keystone CR.
	KeystoneConf *contrail.KeystoneConfiguration
}

type extKeystoneClient struct {
	url    string
	client http.Client
}

func newExtKeystoneClient(protocol string, address string, port int, ca []byte) *extKeystoneClient {
	var httpClient http.Client

	if protocol == "https" && ca != nil {
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
		rootCAs.AppendCertsFromPEM(ca)
		config := &tls.Config{
			RootCAs: rootCAs,
		}
		tr := &http.Transport{TLSClientConfig: config}
		httpClient = http.Client{Transport: tr}
	} else {
		httpClient = http.Client{}
	}

	url := fmt.Sprintf("%s://%s:%d", protocol, address, port)
	return &extKeystoneClient{
		url:    url,
		client: httpClient,
	}
}

func (c *extKeystoneClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.url+path, body)
}

func (c *extKeystoneClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) PostAuthTokens(username, password, project string) (AuthTokens, error) {
	return c.PostAuthTokensWithHeaders(username, password, project, http.Header{})
}

func (c *Client) PostAuthTokensWithHeaders(username, password, project string, headers http.Header) (AuthTokens, error) {
	kar := &keystoneAuthRequest{}
	kar.Auth.Identity.Methods = []string{"password"}
	kar.Auth.Identity.Password.User.Name = username
	kar.Auth.Identity.Password.User.Domain.ID = c.KeystoneConf.UserDomainID
	kar.Auth.Identity.Password.User.Password = password
	kar.Auth.Scope.Project.Domain.ID = c.KeystoneConf.ProjectDomainID
	kar.Auth.Scope.Project.Name = project
	karBody, err := json.Marshal(kar)
	if err != nil {
		return AuthTokens{}, err
	}
	request, err := c.Connector.NewRequest(http.MethodPost, "/v3/auth/tokens", bytes.NewReader(karBody))
	if err != nil {
		return AuthTokens{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	for name, values := range headers {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	response, err := c.Connector.Do(request)
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
