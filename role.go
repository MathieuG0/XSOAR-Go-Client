package xsoar

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MathieuG0/XSOAR-Go-Client/cache"
)

type Shift struct {
	FromDay    int `json:"fromDay"`
	FromHour   int `json:"fromHour"`
	FromMinute int `json:"fromMinute"`
	ToDay      int `json:"toDay"`
	ToHour     int `json:"toHour"`
	ToMinute   int `json:"toMinute"`
}

type PredefinedRange struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Picker struct {
	PredefinedRange   PredefinedRange `json:"predefinedRange"`
	RelativeTimeRange any             `json:"relativeTimeRange"`
	Start             time.Time       `json:"start"`
}

type TableColumn struct {
	IsDefault bool   `json:"isDefault"`
	Key       string `json:"key"`
	Position  int    `json:"position"`
	Width     int    `json:"width"`
}

type UserPreferencesIncidentTableQueries struct {
	ID     string `json:"id"`
	Picker Picker `json:"picker"`
	Query  string `json:"query"`
}

type UserPreferencesIndicatorsTableQueries struct {
	ID           string        `json:"id"`
	Picker       Picker        `json:"picker"`
	Query        string        `json:"query"`
	TableColumns []TableColumn `json:"tableColumn"`
	TableSortAsc bool          `json:"tableSortAsc"`
	TableSortKey string        `json:"tableSortKey"`
}

type UserPreferencesJobsTableQueries struct {
	ID           string `json:"id"`
	Picker       Picker `json:"picker"`
	Query        string `json:"query"`
	TableSortAsc bool   `json:"tableSortAsc"`
	TableSortKey string `json:"tableSortKey"`
}

type UserPreferencesWarRoomFilterMap struct {
	Categories       []string  `json:"categories"`
	FromTime         time.Time `json:"fromTime"`
	ID               string    `json:"id"`
	PageSize         int       `json:"pageSize"`
	Query            string    `json:"query"`
	TagsAndOperator  bool      `json:"tagsAndOperator"`
	UsersAndOperator bool      `json:"usersAndOperator"`
}

type DefaultPreferences struct {
	UserPreferencesIncidentTableQueries   UserPreferencesIncidentTableQueries   `json:"userPreferencesIncidentTableQueries"`
	UserPreferencesIndicatorsTableQueries UserPreferencesIndicatorsTableQueries `json:"userPreferencesIndicatorsTableQueries"`
	UserPreferencesWarRoomFilterMap       UserPreferencesWarRoomFilterMap       `json:"userPreferencesWarRoomFilterMap"`
	UserPreferencesJobsTableQueries       UserPreferencesJobsTableQueries       `json:"userPreferencesJobsTableQueries"`
}

type Role struct {
	ID                 string             `json:"id"`
	Version            int                `json:"version"`
	Name               string             `json:"name"`
	Locked             bool               `json:"locked"`
	NestedRoles        []string           `json:"nestedRoles"`
	AllRoles           []string           `json:"allRoles"`
	Permissions        []string           `json:"permissions"`
	AllPermissions     []string           `json:"allPermissions"`
	ADGroups           []string           `json:"adGroups"`
	SamlGroups         []string           `json:"samlGroups"`
	PagesAccess        []string           `json:"pagesAccess"`
	Shifts             []Shift            `json:"shifts"`
	DefaultPreferences DefaultPreferences `json:"defaultPreferences"`
	DefaultDashboards  []string           `json:"defaultDashboards"`
}

type RoleModule struct {
	client *Client
	cache  *cache.Cache
}

func (m *RoleModule) GetRoles() ([]Role, error) {
	req, err := m.client.NewRequest(
		http.MethodGet, "roles",
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return HTTPResponseDecode[[]Role](resp)
}

func (m *RoleModule) UpsertRole(r Role) ([]Role, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "roles/update",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return HTTPResponseDecode[[]Role](resp)
}

func (m *RoleModule) DeleteRole(id string) ([]Role, error) {
	req, err := m.client.NewRequest(
		http.MethodDelete, "roles/"+id,
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return HTTPResponseDecode[[]Role](resp)
}
