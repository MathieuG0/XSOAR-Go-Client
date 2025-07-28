package xsoar

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type APIKey struct {
	CacheVersn  int       `json:"cacheVersn"`
	Created     time.Time `json:"created"`
	ID          string    `json:"id"`
	Locked      bool      `json:"locked"`
	Modified    time.Time `json:"modified"`
	Name        string    `json:"name"`
	SizeInBytes int       `json:"sizeInBytes"`
	Username    string    `json:"username"`
	Version     int       `json:"version"`
	APIKey      string    `json:"-"`
}

type APIKeyCreate struct {
	Name   string `json:"name"`
	APIKey string `json:"apikey"`
}

func generateKey() string {
	d := make([]byte, 16)
	_, _ = rand.Read(d)
	return strings.ToUpper(hex.EncodeToString(d))
}

func (m *IntegrationModule) ListAPIKeys() ([]APIKey, error) {
	req, err := m.client.NewRequest(
		http.MethodGet, "apikeys",
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return Decode[[]APIKey](resp)
}

func (m *IntegrationModule) CreateAPIKey(name string) (APIKey, error) {
	payload := APIKeyCreate{Name: name, APIKey: generateKey()}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return APIKey{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "apikeys",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return APIKey{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return APIKey{}, err
	}

	keys, err := Decode[[]APIKey](resp)
	if err != nil {
		return APIKey{}, err
	}

	for _, key := range keys {
		if key.Name == payload.Name {
			key.APIKey = payload.APIKey
			return key, nil
		}
	}

	return APIKey{}, errors.Errorf("key not found after creation")
}

func (m *IntegrationModule) DeleteAPIKey(id string) ([]APIKey, error) {
	req, err := m.client.NewRequest(
		http.MethodDelete, "apikeys/"+id,
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return Decode[[]APIKey](resp)
}
