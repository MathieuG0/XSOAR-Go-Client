package xsoar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type IntegrationParamType int

const (
	AuthenticationParamType IntegrationParamType = 9
	BoolParamType           IntegrationParamType = 8
	EncryptedParamType      IntegrationParamType = 4
	LongTextParamType       IntegrationParamType = 12
	MultiSelectParamType    IntegrationParamType = 16
	ShortTextParamType      IntegrationParamType = 0
	SingleSelectParamType   IntegrationParamType = 15
)

type InstanceIntegrationData struct {
	Section         string               `json:"section"`
	Advanced        bool                 `json:"advanced"`
	Display         string               `json:"display"`
	DisplayPassword string               `json:"displayPassword"`
	Name            string               `json:"name"`
	DefaultValue    string               `json:"defaultValue"`
	Type            IntegrationParamType `json:"type"`
	Required        bool                 `json:"required"`
	Hidden          bool                 `json:"hidden"`
	HiddenUsername  bool                 `json:"hiddenUsername"`
	HiddenPassword  bool                 `json:"hiddenPassword"`
	Options         []string             `json:"options"`
	Info            string               `json:"info"`
	Hasvalue        bool                 `json:"hasvalue"`
	Value           json.RawMessage      `json:"value"`
}

type IntegrationPermission struct {
	Roles         []string `json:"roles"`
	PreviousRoles []string `json:"previousRoles"`
	HasRole       bool     `json:"hasRole"`
}

type IntegrationInstance struct {
	ID                     string                           `json:"id"`
	Name                   string                           `json:"name"`
	Brand                  string                           `json:"brand"`
	Version                int                              `json:"version"`
	Category               string                           `json:"category"`
	CacheVersion           int                              `json:"cacheVersn"`
	Modified               time.Time                        `json:"modified"`
	Created                time.Time                        `json:"created"`
	SizeInBytes            int                              `json:"sizeInBytes"`
	SyncHash               string                           `json:"syncHash"`
	PackID                 string                           `json:"packId"`
	PackName               string                           `json:"packName"`
	ItemVersion            string                           `json:"itemVersion"`
	FromServerVersion      string                           `json:"fromServerVersion"`
	ToServerVersion        string                           `json:"ToServerVersion"`
	PropagationLabels      []string                         `json:"PropagationLabels"`
	DefinitionId           string                           `json:"definitionId"`
	PrevName               string                           `json:"prevName"`
	Password               string                           `json:"password"`
	Enabled                bool                             `json:"enabled,string"`
	ConfigValues           json.RawMessage                  `json:"configvalues"`
	ConfigTypes            map[string]IntegrationParamType  `json:"configtypes"`
	Path                   string                           `json:"path"`
	Executable             string                           `json:"executable"`
	Cmdline                string                           `json:"cmdline"`
	Engine                 string                           `json:"engine"`
	EngineGroup            string                           `json:"engineGroup"`
	Hidden                 bool                             `json:"hidden"`
	IsIntegrationScript    bool                             `json:"isIntegrationScript"`
	IsLongRunning          bool                             `json:"islongRunning"`
	MappingId              string                           `json:"mappingId"`
	OutgoingMapperId       string                           `json:"outgoingMapperId"`
	IncomingMapperId       string                           `json:"incomingMapperId"`
	RemoteSync             bool                             `json:"remoteSync"`
	IsSystemIntegration    bool                             `json:"isSystemIntegration"`
	CanSample              bool                             `json:"canSample"`
	DefaultIgnore          bool                             `json:"defaultIgnore"`
	IntegrationLogLevel    string                           `json:"integrationLogLevel"`
	CommandsPermissions    map[string]IntegrationPermission `json:"commandsPermissions"`
	LongRunningId          string                           `json:"longRunningId"`
	IncidentFetchInterval  int                              `json:"incidentFetchInterval"`
	EventFetchInterval     int                              `json:"eventFetchInterval"`
	ServicesID             string                           `json:"servicesID"`
	IsBuiltin              bool                             `json:"isBuiltin"`
	Configuration          Integration                      `json:"configuration"`
	Data                   []InstanceIntegrationData        `json:"data"`
	DisplayPassword        string                           `json:"displayPassword"`
	PasswordProtected      bool                             `json:"passwordProtected"`
	Mappable               bool                             `json:"mappable"`
	RemoteSyncableIn       bool                             `json:"remoteSyncableIn"`
	RemoteSyncableOut      bool                             `json:"remoteSyncableOut"`
	IsFetchSamples         bool                             `json:"isFetchSamples"`
	DebugMode              bool                             `json:"debugMode"`
	UnclassifiedCasesCount int                              `json:"unclassifiedCasesCount"`
}

type IntegrationCommandArgument struct {
	Name         string   `json:"name"`
	Required     bool     `json:"required"`
	Deprecated   bool     `json:"deprecated"`
	Hidden       bool     `json:"hidden"`
	Default      bool     `json:"default"`
	Secret       bool     `json:"secret"`
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
	Type         string   `json:"type"`
	IsArray      bool     `json:"isArray"`
	Auto         string   `json:"auto"`
	Predefined   []string `json:"predefined"`
}

type IntegationCommand struct {
	Name            string                       `json:"name"`
	Deprecated      bool                         `json:"deprecated"`
	Arguments       []IntegrationCommandArgument `json:"arguments"`
	Outputs         any                          `json:"outputs"`
	Important       any                          `json:"important"`
	Description     string                       `json:"description"`
	Execution       bool                         `json:"execution"`
	Cartesian       bool                         `json:"cartesian"`
	Hidden          bool                         `json:"hidden"`
	DocsHidden      bool                         `json:"docsHidden"`
	Sensitive       bool                         `json:"sensitive"`
	Timeout         int                          `json:"timeout"`
	Permitted       bool                         `json:"permitted"`
	IndicatorAction bool                         `json:"indicatorAction"`
	DefinitionId    string                       `json:"definitionId"`
	GomAction       bool                         `json:"gomAction"`
	Polling         bool                         `json:"polling"`
}

type IntegrationScript struct {
	Script                 string              `json:"script"`
	Type                   string              `json:"type"`
	Commands               []IntegationCommand `json:"commands"`
	DockerImage            string              `json:"dockerImage"`
	NativeImage            []string            `json:"nativeImage"`
	IsFetch                bool                `json:"isFetch"`
	IsFetchEvents          bool                `json:"isFetchEvents"`
	IsFetchAssets          bool                `json:"isFetchAssets"`
	Feed                   bool                `json:"feed"`
	IsFetchCredentials     bool                `json:"isFetchCredentials"`
	RunOnce                bool                `json:"runOnce"`
	LongRunning            bool                `json:"longRunning"`
	LongRunningPortMapping bool                `json:"longRunningPortMapping"`
	Subtype                string              `json:"subtype"`
	IsMappable             bool                `json:"isMappable"`
	IsRemoteSyncIn         bool                `json:"isRemoteSyncIn"`
	IsRemoteSyncOut        bool                `json:"isRemoteSyncOut"`
	IsFetchSamples         bool                `json:"isFetchSamples"`
	ResetContext           bool                `json:"resetContext"`
}

type Integration struct {
	ID                                string                    `json:"id"`
	Version                           int                       `json:"version"`
	CacheVersn                        int                       `json:"cacheVersn"`
	SequenceNumber                    int                       `json:"sequenceNumber"`
	PrimaryTerm                       int                       `json:"primaryTerm"`
	Modified                          time.Time                 `json:"modified"`
	Created                           time.Time                 `json:"created"`
	SizeInBytes                       int                       `json:"sizeInBytes"`
	SortValues                        []string                  `json:"sortValues"`
	PackID                            string                    `json:"packID"`
	PackName                          string                    `json:"packName"`
	ItemVersion                       string                    `json:"itemVersion"`
	FromServerVersion                 string                    `json:"fromServerVersion"`
	ToServerVersion                   string                    `json:"toServerVersion"`
	PropagationLabels                 []string                  `json:"propagationLabels"`
	PackPropagationLabels             []string                  `json:"packPropagationLabels"`
	DefinitionId                      string                    `json:"definitionId"`
	VcShouldIgnore                    bool                      `json:"vcShouldIgnore"`
	VcShouldKeepItemLegacyProdMachine bool                      `json:"vcShouldKeepItemLegacyProdMachine"`
	CommitMessage                     string                    `json:"commitMessage"`
	ShouldCommit                      bool                      `json:"shouldCommit"`
	Name                              string                    `json:"name"`
	PrevName                          string                    `json:"prevName"`
	Display                           string                    `json:"display"`
	Brand                             string                    `json:"brand"`
	Category                          string                    `json:"category"`
	Icon                              string                    `json:"icon"`
	Image                             string                    `json:"image"`
	Description                       string                    `json:"description"`
	DetailedDescription               string                    `json:"detailedDescription"`
	SectionOrder                      []string                  `json:"sectionOrder"`
	Configuration                     []InstanceIntegrationData `json:"configuration"`
	ReadOnly                          bool                      `json:"readonly"`
	IntegrationScript                 IntegrationScript         `json:"integrationScript"`
	IsPasswordProtected               bool                      `json:"isPasswordProtected"`
	System                            bool                      `json:"system"`
	Hidden                            bool                      `json:"hidden"`
	CanGetSamples                     bool                      `json:"canGetSamples"`
	Deprecated                        bool                      `json:"deprecated"`
	DefaultMapperIn                   string                    `json:"defaultClassifier"`
	DefaultClassifier                 string                    `json:"defaultMapperIn"`
}

type Engines struct {
	EngineGroups  any      `json:"engineGroups"`
	Engines       []any    `json:"engines"`
	PkgTypes      []string `json:"pkgTypes"`
	RequestedLogs any      `json:"requestedLogs"`
	Total         int      `json:"total"`
}

type InstanceHealth struct {
	ID                    string    `json:"id"`
	Version               int       `json:"version"`
	CacheVersn            int       `json:"cacheVersn"`
	Modified              time.Time `json:"modified"`
	Created               time.Time `json:"created"`
	SizeInBytes           int       `json:"sizeInBytes"`
	Brand                 string    `json:"brand"`
	Instance              string    `json:"instance"`
	IncidentsPulled       int       `json:"incidentsPulled"`
	IndicatorsPulled      int       `json:"indicatorsPulled"`
	EventsPulled          int       `json:"eventsPulled"`
	IncidentsDropped      int       `json:"incidentsDropped"`
	LastPullTime          time.Time `json:"lastPullTime"`
	LastError             string    `json:"lastError"`
	LastIngestionDuration int       `json:"lastIngestionDuration"`
	FetchDuration         int       `json:"fetchDuration"`
	Creation_time         time.Time `json:"creation_time"`
	Sum_hour              int       `json:"sum_hour"`
	Sum_day               int       `json:"sum_day"`
	Sum_week              int       `json:"sum_week"`
}

type IntegrationSearch struct {
	Configurations []Integration             `json:"configurations"`
	Engines        Engines                   `json:"engines"`
	Health         map[string]InstanceHealth `json:"health"`
	Instances      []IntegrationInstance     `json:"instances"`
}

type IntegrationCommands struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Display             string              `json:"display"`
	Category            string              `json:"category"`
	Description         string              `json:"description"`
	DetailedDescription string              `json:"detailedDescription"`
	Commands            []IntegationCommand `json:"commands"`
	Feed                bool                `json:"feed"`
	IsFetch             bool                `json:"isFetch"`
}

type InstanceIntegrationDataUpsert struct {
	Name     string               `json:"name"`
	Type     IntegrationParamType `json:"type"`
	Value    any                  `json:"value"`
	Hasvalue bool                 `json:"hasvalue"`
}

type IntegrationInstanceUpsert struct {
	ID                  string                          `json:"id,omitempty"`
	Name                string                          `json:"name"`
	Brand               string                          `json:"brand"`
	Version             int                             `json:"version"`
	Enabled             bool                            `json:"enabled,string"`
	Engine              string                          `json:"engine"`
	EngineGroup         string                          `json:"engineGroup"`
	Hidden              bool                            `json:"hidden"`
	IsIntegrationScript bool                            `json:"isIntegrationScript"`
	MappingId           string                          `json:"mappingId"`
	OutgoingMapperId    string                          `json:"outgoingMapperId"`
	IncomingMapperId    string                          `json:"incomingMapperId"`
	CanSample           bool                            `json:"canSample"`
	IntegrationLogLevel string                          `json:"integrationLogLevel"`
	PropagationLabels   []string                        `json:"PropagationLabels"`
	DefaultIgnore       bool                            `json:"defaultIgnore"`
	Data                []InstanceIntegrationDataUpsert `json:"data,omitempty"`
	// CommandsPermissions   map[string]IntegrationPermission `json:"commandsPermissions"`
}

type SearchIntegrationsOptions struct {
	InstanceID string
}

type IntegrationModule struct {
	client *Client
}

func (m *IntegrationModule) GetInstances() ([]IntegrationInstance, error) {
	req, err := m.client.NewRequest(http.MethodGet, "integration/instances", WithHeader("Accept", "application/json"))
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return Decode[[]IntegrationInstance](resp)
}

func (m *IntegrationModule) UpsertInstance(instance IntegrationInstanceUpsert) (IntegrationInstance, error) {
	instance.IsIntegrationScript = true

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(instance); err != nil {
		return IntegrationInstance{}, err
	}

	req, err := m.client.NewRequest(
		http.MethodPut, "settings/integration",
		WithBody(buf),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return IntegrationInstance{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return IntegrationInstance{}, err
	}

	return Decode[IntegrationInstance](resp)
}

func (m *IntegrationModule) DeleteInstance(id string) error {
	req, err := m.client.NewRequest(
		http.MethodDelete, fmt.Sprintf("settings/integration/%s", id),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return err
	}

	_, err = m.client.Do(req)
	return err
}

func (m *IntegrationModule) SearchIntegrations(opt *SearchIntegrationsOptions) (IntegrationSearch, error) {
	req, err := m.client.NewRequest(
		http.MethodPost, "settings/integration/search",
		WithBody(strings.NewReader(`{}`)),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Accept", "application/json"),
	)
	if err != nil {
		return IntegrationSearch{}, err
	}

	if opt != nil && opt.InstanceID != "" {
		q := req.URL.Query()
		q.Set("id", opt.InstanceID)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return IntegrationSearch{}, err
	}

	return Decode[IntegrationSearch](resp)
}

func (m *IntegrationModule) GetIntegrationCommands() ([]IntegrationCommands, error) {
	req, err := m.client.NewRequest(http.MethodGet, "settings/integration-commands", WithHeader("Accept", "application/json"))
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return Decode[[]IntegrationCommands](resp)
}
