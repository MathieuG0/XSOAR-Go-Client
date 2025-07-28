package xsoar

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Credential struct {
	CacheVersn         int       `json:"cacheVersn"`
	Created            time.Time `json:"created"`
	HasCertificate     bool      `json:"hasCertificate"`
	HasCertificatePass bool      `json:"hasCertificatePass"`
	HasPassword        bool      `json:"hasPassword"`
	ID                 string    `json:"id"`
	Locked             bool      `json:"locked"`
	Modified           time.Time `json:"modified"`
	Name               string    `json:"name"`
	SizeInBytes        int       `json:"sizeInBytes"`
	User               string    `json:"user"`
	VaultInstanceId    string    `json:"vaultInstanceId"`
	Version            int       `json:"version"`
	Workgroup          string    `json:"workgroup"`
}

type CredentialUpsert struct {
	ID             string `json:"id,omitempty"`
	HasCertificate bool   `json:"hasCertificate,omitempty"`
	HasPassword    bool   `json:"hasCertificatePass,omitempty"`
	Name           string `json:"name,omitempty"`
	Password       string `json:"password,omitempty"`
	SSHKey         string `json:"sshkey,omitempty"`
	User           string `json:"user,omitempty"`
	Version        int    `json:"version,omitempty"`
	Workgroup      string `json:"workgroup,omitempty"`
}

type CredentialSearch struct {
	Credentials []Credential `json:"credentials"`
	Total       int          `json:"total"`
}

type CredentialDelete struct {
	ID string `json:"id"`
}

func (m *IntegrationModule) ListCredentials() (CredentialSearch, error) {
	req, err := m.client.NewRequest(
		http.MethodPost, "settings/credentials",
		WithBody(strings.NewReader(`{"page":0,"size":200,"query":""}`)),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return CredentialSearch{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return CredentialSearch{}, err
	}

	return Decode[CredentialSearch](resp)
}

func (m *IntegrationModule) UpsertCredential(credential CredentialUpsert) (Credential, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(credential); err != nil {
		return Credential{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPut, "settings/credentials",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return Credential{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return Credential{}, err
	}

	return Decode[Credential](resp)
}

func (m *IntegrationModule) DeleteCredential(id string) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(CredentialDelete{id}); err != nil {
		return err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "settings/credentials/delete",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return err
	}

	_, err = m.client.Do(req)
	return err
}
