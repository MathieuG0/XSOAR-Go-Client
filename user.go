package xsoar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MathieuG0/XSOAR-Go-Client/cache"
)

type User struct {
	LastLogin             time.Time           `json:"lastLogin"`
	LastLoginMaster       time.Time           `json:"lastLoginMaster"`
	Roles                 map[string][]string `json:"roles"`
	AllRoles              []string            `json:"allRoles"`
	Permissions           map[string][]string `json:"permissions"`
	PagesAccess           []string            `json:"PagesAccess"`
	ID                    string              `json:"id"`
	Username              string              `json:"username"`
	Email                 string              `json:"email"`
	Phone                 string              `json:"phone"`
	FirstName             string              `json:"firstName"`
	LastName              string              `json:"lastName"`
	Name                  string              `json:"name"`
	UserType              string              `json:"userType"`
	PlaygroundId          string              `json:"playgroundId"`
	DefaultAdmin          bool                `json:"defaultAdmin"`
	AccUser               bool                `json:"accUser"`
	IsLocked              bool                `json:"isLocked"`
	IsAway                bool                `json:"isAway"`
	PlaygroundCleared     bool                `json:"playgroundCleared"`
	Disabled              bool                `json:"disabled"`
	ReadOnly              bool                `json:"readOnly"`
	Homepage              string              `json:"homepage"`
	InvestigationPage     string              `json:"investigationPage"`
	EditorStyle           string              `json:"editorStyle"`
	TimeZone              string              `json:"timeZone"`
	UserTimeZone          string              `json:"userTimeZone"`
	DateFormat            string              `json:"dateFormat"`
	TimeFormat            string              `json:"timeFormat"`
	Hours24               string              `json:"hours24"`
	Theme                 string              `json:"theme"`
	Image                 string              `json:"image"`
	WasAssigned           bool                `json:"wasAssigned"`
	Accounts              []string            `json:"accounts"`
	DisableHyperSearch    bool                `json:"disableHyperSearch"`
	HelpSnippetDisabled   bool                `json:"helpSnippetDisabled"`
	ShortcutsDisabled     bool                `json:"shortcutsDisabled"`
	Dashboards            json.RawMessage     `json:"dashboards"`
	AddedSharedDashboards []string            `json:"addedSharedDashboards"`
	Type                  int                 `json:"type"`
	Preferences           json.RawMessage     `json:"preferences"`
	NotificationsSettings json.RawMessage     `json:"notificationsSettings"`
}

type InviteCreation struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type Invite struct {
	ID            string    `json:"id"`
	Version       int       `json:"version"`
	CacheVersn    int       `json:"cacheVersn"`
	Modified      time.Time `json:"modified"`
	Created       time.Time `json:"created"`
	SizeInBytes   int       `json:"sizeInBytes"`
	CreatedBy     string    `json:"createdBy"`
	Email         string    `json:"email"`
	Expiration    time.Time `json:"expiration"`
	Roles         []string  `json:"roles"`
	Url           string    `json:"url"`
	Project       string    `json:"project"`
	Accepted      int       `json:"accepted"`
	Investigation string    `json:"investigation"`
}

type InviteUtilization struct {
	ID       string `json:"-"`
	Existing bool   `json:"existing"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserPasswordReset struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type UserRoleUpdate struct {
	ID    string              `json:"id"`
	Roles UserRoleUpdateRoles `json:"roles"`
}

type UserRoleUpdateRoles struct {
	DefaultAdmin bool     `json:"defaultAdmin"`
	Roles        []string `json:"roles"`
}

type InviteSearch struct {
	Total   int      `json:"total"`
	Invites []Invite `json:"invites"`
}

type UserModule struct {
	client *Client
	cache  *cache.Cache
}

func (m *UserModule) GetUsers() ([]User, error) {
	req, err := m.client.NewRequest(
		http.MethodGet, "users",
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return HTTPResponseDecode[[]User](resp)
}

func (m *UserModule) CreateInvite(i InviteCreation) (Invite, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(i); err != nil {
		return Invite{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "invite",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return Invite{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return Invite{}, err
	}

	return HTTPResponseDecode[Invite](resp)
}

func (m *UserModule) UtilizeInvite(i InviteUtilization) (User, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(i); err != nil {
		return User{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, fmt.Sprintf("invite/%s/utilize", i.ID),
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return User{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return User{}, err
	}

	return HTTPResponseDecode[User](resp)
}

func (m *UserModule) DeleteInvite(ids ...string) (InviteSearch, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(map[string][]string{"ids": ids}); err != nil {
		return InviteSearch{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "invites/delete",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return InviteSearch{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return InviteSearch{}, err
	}

	return HTTPResponseDecode[InviteSearch](resp)
}

func (m *UserModule) ResetPassword(p UserPasswordReset) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		return err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "users/setpw",
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

func (m *UserModule) Disable(id string) ([]User, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(map[string]string{"id": id}); err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "users/disable",
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

	return HTTPResponseDecode[[]User](resp)
}

func (m *UserModule) Enable(id string) ([]User, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(map[string]string{"id": id}); err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "users/enable",
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

	return HTTPResponseDecode[[]User](resp)
}

func (m *UserModule) Update(u UserRoleUpdate) ([]User, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(u); err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "users/update",
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

	return HTTPResponseDecode[[]User](resp)
}

func (m *UserModule) Delete(ids ...string) ([]User, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(map[string][]string{"ids": ids}); err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest(
		http.MethodPost, "users/delete",
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

	return HTTPResponseDecode[[]User](resp)
}
