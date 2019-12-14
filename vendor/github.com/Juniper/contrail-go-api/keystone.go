//
// Copyright (c) 2014 Juniper Networks, Inc. All rights reserved.
//

package contrail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
)

// KeystoneClient is a client of the OpenStack Keystone service that adds authentication
// tokens to the Contrail API requests.
type KeystoneClient struct {
	osAuthURL    string
	osTenantName string
	osUsername   string
	osPassword   string
	osAdminToken string

	current *KeystoneToken
}

type KeepaliveKeystoneClient struct {
	KeystoneClient
}

// KeystoneToken represents an auth token issued by OpenStack keystone service.
// The field names are defined by the Keystone API schema.
type KeystoneToken struct {
	Id      string
	Expires string
	Tenant  struct {
		Id          string
		Name        string
		Description string
		Enabled     bool
	}
	Issued_At string
}

// NewKeystoneClient allocates and initializes a KeystoneClient
func NewKeystoneClient(auth_url, tenant_name, username, password, token string) *KeystoneClient {
	return &KeystoneClient{
		auth_url,
		tenant_name,
		username,
		password,
		token,
		nil,
	}
}

func NewKeepaliveKeystoneClient(auth_url, tenant_name, username, password, token string) *KeepaliveKeystoneClient {
	return &KeepaliveKeystoneClient {
		KeystoneClient{
			auth_url,
			tenant_name,
			username,
			password,
			token,
			nil,
		},
	}
}

// Authenticate sends an authentication request to keystone.
func (kClient *KeystoneClient) Authenticate() error {
	// identity:CredentialType
	type AuthTokenRequest struct {
		Auth struct {
			Token struct {
				Id string `json:"id"`
			} `json:"token"`
		} `json:"auth"`
	}
	type AuthCredentialsRequest struct {
		Auth struct {
			TenantName          string `json:"tenantName"`
			PasswordCredentials struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"passwordCredentials"`
		} `json:"auth"`
	}
	// identity-api/v2.0/src/xsd/token.xsd
	// <element name="access" type="identity:AuthenticateResponse"/>
	type TokenResponse struct {
		Access struct {
			Token KeystoneToken
			User  struct {
				Id       string
				Username string
			}
			// ServiceCatalog
		}
	}
	url := kClient.osAuthURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "tokens"

	var data []byte
	var err error
	if len(kClient.osAdminToken) > 0 {
		request := AuthTokenRequest{}
		request.Auth.Token.Id = kClient.osAdminToken
		data, err = json.Marshal(&request)
	} else {
		request := AuthCredentialsRequest{}
		request.Auth.PasswordCredentials.Username =
			kClient.osUsername
		request.Auth.PasswordCredentials.Password =
			kClient.osPassword
		request.Auth.TenantName = kClient.osTenantName
		data, err = json.Marshal(&request)
	}
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json",
		bytes.NewReader(data))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", resp.Status, body)
	}

	var response TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	kClient.current = new(KeystoneToken)
	*kClient.current = response.Access.Token
	return nil
}

func (kClient *KeepaliveKeystoneClient) needsRefreshing() (bool, error) {
	if kClient.current == nil {
		return true, nil
	}

	issuedAtTime, err := time.Parse(time.RFC3339, kClient.current.Issued_At)
	if err != nil {
		return false, err
	}

	expires, err := time.Parse(time.RFC3339, kClient.current.Expires)
	if err != nil {
		return false, err
	}

	refreshTime := issuedAtTime.UTC().Add(expires.UTC().Sub(issuedAtTime.UTC()) / 2)

	return time.Now().UTC().After(refreshTime.UTC()), nil
}

func (kClient *KeepaliveKeystoneClient) AddAuthentication(req *http.Request) error {
	needsRefreshing, err := kClient.needsRefreshing()
	if err != nil {
		return err
	}

	if needsRefreshing {
		kClient.current = nil
	}

	return kClient.KeystoneClient.AddAuthentication(req)
}

// AddAuthentication adds the authentication data to the HTTP header.
func (kClient *KeystoneClient) AddAuthentication(req *http.Request) error {
	if kClient.current == nil {
		err := kClient.Authenticate()
		if err != nil {
			return err
		}
	}
	req.Header.Set("X-Auth-Token", kClient.current.Id)
	return nil
}
