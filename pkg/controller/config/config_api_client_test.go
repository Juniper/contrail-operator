package config_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"atom/atom/contrail/operator/pkg/controller/config"
)

func TestEnsureConfigNodeExists(t *testing.T) {

	t.Run("should create config API client", func(t *testing.T) {
		client := config.NewApiClient("localhost:8082")
		assert.NotNil(t, client)
	})

	t.Run("should POST /config-nodes", func(t *testing.T) {
		tests := map[string]struct {
			configNode   config.ConfigNode
			expectedBody string
		}{
			"localhost": {
				configNode: config.ConfigNode{
					UUID:     "520f1126-34cc-4d1f-bda8-4df1b5aeea7d",
					Hostname: "localhost",
					IP:       "10.0.2.15",
				},
				expectedBody: `
					{ 
					   "config-node":{ 
						  "parent_type":"global-system-config",
						  "fq_name":[ 
							 "default-global-system-config",
							 "localhost"
						  ],
						  "uuid":"520f1126-34cc-4d1f-bda8-4df1b5aeea7d",
						  "config_node_ip_address":"10.0.2.15"
					   }
					}`,
			},
			"juniper.net": {
				configNode: config.ConfigNode{
					UUID:     "86f38811-a892-4877-885f-be0fa05ea164",
					Hostname: "juniper.net",
					IP:       "10.0.2.1",
				},
				expectedBody: `
					{ 
					   "config-node":{ 
						  "parent_type":"global-system-config",
						  "fq_name":[ 
							 "default-global-system-config",
							 "juniper.net"
						  ],
						  "uuid":"86f38811-a892-4877-885f-be0fa05ea164",
						  "config_node_ip_address":"10.0.2.1"
					   }
					}`,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// given
				server := startHTTPServerStub(200)
				defer server.close()
				client := config.NewApiClient(server.url())
				// when
				err := client.EnsureConfigNodeExists(test.configNode)
				req := server.waitForRequest()
				// then
				assert.NoError(t, err)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/config-nodes", req.RequestURI)
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
				assert.JSONEq(t, test.expectedBody, req.body)
			})
		}
	})

	t.Run("should return err when server is down", func(t *testing.T) {
		// given
		client := config.NewApiClient("http://127.0.0.1:1")
		// when
		err := client.EnsureConfigNodeExists(config.ConfigNode{
			UUID:     "123-123-123-123",
			Hostname: "localhost",
		})
		// then
		assert.Error(t, err)
	})

	t.Run("should return err when server returned http error", func(t *testing.T) {
		statusCodes := []int{404, 400, 500}
		for _, statusCode := range statusCodes {
			t.Run(strconv.Itoa(statusCode), func(t *testing.T) {
				// given
				server := startHTTPServerStub(statusCode)
				defer server.close()
				client := config.NewApiClient(server.url())
				// when
				err := client.EnsureConfigNodeExists(config.ConfigNode{
					UUID:     "123-123-123-123",
					Hostname: "localhost",
				})
				// then
				assert.Error(t, err)
			})
		}
	})

	t.Run("when config node already exists should return without error", func(t *testing.T) {
		// given
		server := startHTTPServerStub(409)
		defer server.close()
		client := config.NewApiClient(server.url())
		// when
		err := client.EnsureConfigNodeExists(config.ConfigNode{
			UUID:     "123-123-123-123",
			Hostname: "localhost",
		})
		// then
		assert.NoError(t, err)
	})

}

type httpServerStub struct {
	requestReceived chan *request
	server          *httptest.Server
}

func startHTTPServerStub(statusCode int) *httpServerStub {
	stub := &httpServerStub{
		requestReceived: make(chan *request, 1),
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(statusCode)
		stub.requestReceived <- newRequestReceived(req)
	}))
	stub.server = server
	return stub
}

func (s *httpServerStub) waitForRequest() *request {
	return <-s.requestReceived
}

func (s *httpServerStub) url() string {
	return s.server.URL
}

func (s *httpServerStub) close() {
	s.server.Close()
}

type request struct {
	*http.Request
	body string
}

func newRequestReceived(req *http.Request) *request {
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	return &request{Request: req, body: string(body)}
}
