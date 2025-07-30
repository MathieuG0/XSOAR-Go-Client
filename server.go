package xsoar

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/MathieuG0/XSOAR-Go-Client/cache"
	"github.com/pkg/errors"
)

type SystemConfigDefaults struct {
	HTTPProxy        string `json:"http_proxy"`
	HTTPSProxy       string `json:"https_proxy"`
	BaseURL          string `json:"server.baseurl"`
	ExternalHostname string `json:"server.externalhostname"`
}

type systemConfig struct {
	DefaultMap SystemConfigDefaults `json:"defaultMap"`
	SysConfig  map[string]any       `json:"sysConf"`
}

type SystemConfig struct {
	DefaultMap SystemConfigDefaults
	SysConfig  map[string]string
	Version    int
}

type SystemConfigUpdate struct {
	Data    map[string]string `json:"data"`
	Version int               `json:"version"`
}

type ServerModule struct {
	client *Client
	cache  *cache.Cache
}

func toSystemConfig(input systemConfig) (SystemConfig, error) {
	result := SystemConfig{
		DefaultMap: input.DefaultMap,
		SysConfig:  make(map[string]string),
	}

	for key, value := range input.SysConfig {
		if key == "versn" {
			version, ok := value.(float64)
			if !ok {
				return SystemConfig{}, errors.Errorf("invalid version type, expected float64 got: %s", reflect.TypeOf(value))
			}
			result.Version = int(version)
			continue
		}
		str, ok := value.(string)
		if !ok {
			return SystemConfig{}, errors.Errorf("invalid sysconfig value type, expected string got: %s", reflect.TypeOf(value))
		}
		result.SysConfig[key] = str
	}

	return result, nil
}

func (m *ServerModule) GetConfig() (SystemConfig, error) {
	req, err := m.client.NewRequest(
		http.MethodGet, "system/config",
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return SystemConfig{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return SystemConfig{}, err
	}

	config, err := HTTPResponseDecode[systemConfig](resp)
	if err != nil {
		return SystemConfig{}, err
	}

	return toSystemConfig(config)
}

func (m *ServerModule) UpdateConfig(c SystemConfigUpdate) (SystemConfig, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(c); err != nil {
		return SystemConfig{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "system/config",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return SystemConfig{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return SystemConfig{}, err
	}

	config, err := HTTPResponseDecode[systemConfig](resp)
	if err != nil {
		return SystemConfig{}, err
	}

	return toSystemConfig(config)
}
