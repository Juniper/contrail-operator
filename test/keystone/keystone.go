package keystone

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"github.com/Juniper/contrail-operator/test/kubeproxy"
)

func NewClient(t *testing.T, client *kubeproxy.Client) *Client {
	return &Client{client: client, t: t}
}

type Client struct {
	client *kubeproxy.Client
	t      *testing.T
}

func (c *Client) GetAuthTokens(username, password string) AuthTokens {
	kar := &keystoneAuthRequest{}
	kar.Auth.Identity.Methods = []string{"password"}
	kar.Auth.Identity.Password.User.Name = username
	kar.Auth.Identity.Password.User.Domain.ID = "default"
	kar.Auth.Identity.Password.User.Password = password
	karBody, err := json.Marshal(kar)
	require.NoError(c.t, err)
	request := c.client.NewRequest(http.MethodPost, "/v3/auth/tokens", bytes.NewReader(karBody))
	request.Header.Set("Content-Type", "application/json")
	response := c.client.Do(request)
	assert.Equal(c.t, 201, response.StatusCode)
	authResponse := AuthTokens{}
	bytesRead, err := ioutil.ReadAll(response.Body)
	require.NoError(c.t, err)
	err = json.Unmarshal(bytesRead, &authResponse)
	require.NoError(c.t, err)
	authResponse.XAuthTokenHeader = response.Header.Get("X-Subject-Token")
	return authResponse
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

func (t AuthTokens) GetEndpointURL(serviceName string, endpointInterface string) string {
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
