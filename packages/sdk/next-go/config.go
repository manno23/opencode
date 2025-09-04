// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"
	"reflect"

	"github.com/sst/opencode-sdk-go/internal/apijson"
	"github.com/sst/opencode-sdk-go/internal/apiquery"
	"github.com/sst/opencode-sdk-go/internal/param"
	"github.com/sst/opencode-sdk-go/internal/requestconfig"
	"github.com/sst/opencode-sdk-go/option"
	"github.com/sst/opencode-sdk-go/shared"
	"github.com/tidwall/gjson"
)

// ConfigService contains methods and other services that help with interacting
// with the opencode API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewConfigService] method instead.
type ConfigService struct {
	Options []option.RequestOption
}

// NewConfigService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewConfigService(opts ...option.RequestOption) (r *ConfigService) {
	r = &ConfigService{}
	r.Options = opts
	return
}

// Get config info
func (r *ConfigService) Get(ctx context.Context, query ConfigGetParams, opts ...option.RequestOption) (res *Config, err error) {
	opts = append(r.Options[:], opts...)
	path := "config"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

type Config struct {
	// JSON schema reference for configuration validation
	Schema string
	// Agent configuration, see https://opencode.ai/docs/agent
	Agent ConfigAgent
	// @deprecated Use 'share' field instead. Share newly created sessions
	// automatically
	Autoshare bool
	// Automatically update to the latest version
	Autoupdate bool
	// Command configuration, see https://opencode.ai/docs/commands
	Command map[string]ConfigCommand
	// Disable providers that are loaded automatically
	DisabledProviders []string
	Experimental      ConfigExperimental
	Formatter         map[string]ConfigFormatter
	// Additional instruction files or patterns to include
	Instructions []string
	// Custom keybind configurations
	Keybinds KeybindsConfig
	// @deprecated Always uses stretch layout.
	Layout ConfigLayout
	Lsp    map[string]ConfigLsp
	// MCP (Model Context Protocol) server configurations
	Mcp map[string]ConfigMcp
	// @deprecated Use `agent` field instead.
	Mode ConfigMode
	// Model to use in the format of provider/model, eg anthropic/claude-2
	Model      string
	Permission ConfigPermission
	Plugin     []string
	// Custom provider configurations and model overrides
	Provider map[string]ConfigProvider
	// Control sharing behavior:'manual' allows manual sharing via commands, 'auto'
	// enables automatic sharing, 'disabled' disables all sharing
	Share ConfigShare
	// Small model to use for tasks like title generation in the format of
	// provider/model
	SmallModel string
	Snapshot   bool
	// Theme name to use for the interface
	Theme string
	Tools map[string]bool
	// TUI specific settings
	Tui ConfigTui
	// Custom username to display in conversations instead of system username
	Username string
	JSON     configJSON
}

// configJSON contains the JSON metadata for the struct [Config]
type configJSON struct {
	Schema            apijson.Field
	Agent             apijson.Field
	Autoshare         apijson.Field
	Autoupdate        apijson.Field
	Command           apijson.Field
	DisabledProviders apijson.Field
	Experimental      apijson.Field
	Formatter         apijson.Field
	Instructions      apijson.Field
	Keybinds          apijson.Field
	Layout            apijson.Field
	Lsp               apijson.Field
	Mcp               apijson.Field
	Mode              apijson.Field
	Model             apijson.Field
	Permission        apijson.Field
	Plugin            apijson.Field
	Provider          apijson.Field
	Share             apijson.Field
	SmallModel        apijson.Field
	Snapshot          apijson.Field
	Theme             apijson.Field
	Tools             apijson.Field
	Tui               apijson.Field
	Username          apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *Config) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configJSON) RawJSON() string {
	return r.raw
}

// Agent configuration, see https://opencode.ai/docs/agent
type ConfigAgent struct {
	Build       ConfigAgentBuild
	General     ConfigAgentGeneral
	Plan        ConfigAgentPlan
	ExtraFields map[string]ConfigAgent
	JSON        configAgentJSON
}

// configAgentJSON contains the JSON metadata for the struct [ConfigAgent]
type configAgentJSON struct {
	Build       apijson.Field
	General     apijson.Field
	Plan        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentJSON) RawJSON() string {
	return r.raw
}

type ConfigAgentBuild struct {
	// Description of when to use the agent
	Description string
	Disable     bool
	Mode        ConfigAgentBuildMode
	Model       string
	Permission  ConfigAgentBuildPermission
	Prompt      string
	Temperature float64
	Tools       map[string]bool
	TopP        float64
	ExtraFields map[string]interface{}
	JSON        configAgentBuildJSON
}

// configAgentBuildJSON contains the JSON metadata for the struct
// [ConfigAgentBuild]
type configAgentBuildJSON struct {
	Description apijson.Field
	Disable     apijson.Field
	Mode        apijson.Field
	Model       apijson.Field
	Permission  apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	Tools       apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgentBuild) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentBuildJSON) RawJSON() string {
	return r.raw
}

type ConfigAgentBuildMode string

const (
	ConfigAgentBuildModeSubagent ConfigAgentBuildMode = "subagent"
	ConfigAgentBuildModePrimary  ConfigAgentBuildMode = "primary"
	ConfigAgentBuildModeAll      ConfigAgentBuildMode = "all"
)

func (r ConfigAgentBuildMode) IsKnown() bool {
	switch r {
	case ConfigAgentBuildModeSubagent, ConfigAgentBuildModePrimary, ConfigAgentBuildModeAll:
		return true
	}
	return false
}

type ConfigAgentBuildPermission struct {
	Bash     ConfigAgentBuildPermissionBashUnion
	Edit     ConfigAgentBuildPermissionEdit
	Webfetch ConfigAgentBuildPermissionWebfetch
	JSON     configAgentBuildPermissionJSON
}

// configAgentBuildPermissionJSON contains the JSON metadata for the struct
// [ConfigAgentBuildPermission]
type configAgentBuildPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgentBuildPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentBuildPermissionJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ConfigAgentBuildPermissionBashString] or
// [ConfigAgentBuildPermissionBashMap].
type ConfigAgentBuildPermissionBashUnion interface {
	implementsConfigAgentBuildPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigAgentBuildPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigAgentBuildPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigAgentBuildPermissionBashMap{}),
		},
	)
}

type ConfigAgentBuildPermissionBashString string

const (
	ConfigAgentBuildPermissionBashStringAsk   ConfigAgentBuildPermissionBashString = "ask"
	ConfigAgentBuildPermissionBashStringAllow ConfigAgentBuildPermissionBashString = "allow"
	ConfigAgentBuildPermissionBashStringDeny  ConfigAgentBuildPermissionBashString = "deny"
)

func (r ConfigAgentBuildPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionBashStringAsk, ConfigAgentBuildPermissionBashStringAllow, ConfigAgentBuildPermissionBashStringDeny:
		return true
	}
	return false
}

func (r ConfigAgentBuildPermissionBashString) implementsConfigAgentBuildPermissionBashUnion() {}

type ConfigAgentBuildPermissionBashMap map[string]ConfigAgentBuildPermissionBashMapItem

func (r ConfigAgentBuildPermissionBashMap) implementsConfigAgentBuildPermissionBashUnion() {}

type ConfigAgentBuildPermissionBashMapItem string

const (
	ConfigAgentBuildPermissionBashMapAsk   ConfigAgentBuildPermissionBashMapItem = "ask"
	ConfigAgentBuildPermissionBashMapAllow ConfigAgentBuildPermissionBashMapItem = "allow"
	ConfigAgentBuildPermissionBashMapDeny  ConfigAgentBuildPermissionBashMapItem = "deny"
)

func (r ConfigAgentBuildPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionBashMapAsk, ConfigAgentBuildPermissionBashMapAllow, ConfigAgentBuildPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigAgentBuildPermissionEdit string

const (
	ConfigAgentBuildPermissionEditAsk   ConfigAgentBuildPermissionEdit = "ask"
	ConfigAgentBuildPermissionEditAllow ConfigAgentBuildPermissionEdit = "allow"
	ConfigAgentBuildPermissionEditDeny  ConfigAgentBuildPermissionEdit = "deny"
)

func (r ConfigAgentBuildPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionEditAsk, ConfigAgentBuildPermissionEditAllow, ConfigAgentBuildPermissionEditDeny:
		return true
	}
	return false
}

type ConfigAgentBuildPermissionWebfetch string

const (
	ConfigAgentBuildPermissionWebfetchAsk   ConfigAgentBuildPermissionWebfetch = "ask"
	ConfigAgentBuildPermissionWebfetchAllow ConfigAgentBuildPermissionWebfetch = "allow"
	ConfigAgentBuildPermissionWebfetchDeny  ConfigAgentBuildPermissionWebfetch = "deny"
)

func (r ConfigAgentBuildPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionWebfetchAsk, ConfigAgentBuildPermissionWebfetchAllow, ConfigAgentBuildPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigAgentGeneral struct {
	// Description of when to use the agent
	Description string
	Disable     bool
	Mode        ConfigAgentGeneralMode
	Model       string
	Permission  ConfigAgentGeneralPermission
	Prompt      string
	Temperature float64
	Tools       map[string]bool
	TopP        float64
	ExtraFields map[string]interface{}
	JSON        configAgentGeneralJSON
}

// configAgentGeneralJSON contains the JSON metadata for the struct
// [ConfigAgentGeneral]
type configAgentGeneralJSON struct {
	Description apijson.Field
	Disable     apijson.Field
	Mode        apijson.Field
	Model       apijson.Field
	Permission  apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	Tools       apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgentGeneral) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentGeneralJSON) RawJSON() string {
	return r.raw
}

type ConfigAgentGeneralMode string

const (
	ConfigAgentGeneralModeSubagent ConfigAgentGeneralMode = "subagent"
	ConfigAgentGeneralModePrimary  ConfigAgentGeneralMode = "primary"
	ConfigAgentGeneralModeAll      ConfigAgentGeneralMode = "all"
)

func (r ConfigAgentGeneralMode) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralModeSubagent, ConfigAgentGeneralModePrimary, ConfigAgentGeneralModeAll:
		return true
	}
	return false
}

type ConfigAgentGeneralPermission struct {
	Bash     ConfigAgentGeneralPermissionBashUnion
	Edit     ConfigAgentGeneralPermissionEdit
	Webfetch ConfigAgentGeneralPermissionWebfetch
	JSON     configAgentGeneralPermissionJSON
}

// configAgentGeneralPermissionJSON contains the JSON metadata for the struct
// [ConfigAgentGeneralPermission]
type configAgentGeneralPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgentGeneralPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentGeneralPermissionJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ConfigAgentGeneralPermissionBashString] or
// [ConfigAgentGeneralPermissionBashMap].
type ConfigAgentGeneralPermissionBashUnion interface {
	implementsConfigAgentGeneralPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigAgentGeneralPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigAgentGeneralPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigAgentGeneralPermissionBashMap{}),
		},
	)
}

type ConfigAgentGeneralPermissionBashString string

const (
	ConfigAgentGeneralPermissionBashStringAsk   ConfigAgentGeneralPermissionBashString = "ask"
	ConfigAgentGeneralPermissionBashStringAllow ConfigAgentGeneralPermissionBashString = "allow"
	ConfigAgentGeneralPermissionBashStringDeny  ConfigAgentGeneralPermissionBashString = "deny"
)

func (r ConfigAgentGeneralPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionBashStringAsk, ConfigAgentGeneralPermissionBashStringAllow, ConfigAgentGeneralPermissionBashStringDeny:
		return true
	}
	return false
}

func (r ConfigAgentGeneralPermissionBashString) implementsConfigAgentGeneralPermissionBashUnion() {}

type ConfigAgentGeneralPermissionBashMap map[string]ConfigAgentGeneralPermissionBashMapItem

func (r ConfigAgentGeneralPermissionBashMap) implementsConfigAgentGeneralPermissionBashUnion() {}

type ConfigAgentGeneralPermissionBashMapItem string

const (
	ConfigAgentGeneralPermissionBashMapAsk   ConfigAgentGeneralPermissionBashMapItem = "ask"
	ConfigAgentGeneralPermissionBashMapAllow ConfigAgentGeneralPermissionBashMapItem = "allow"
	ConfigAgentGeneralPermissionBashMapDeny  ConfigAgentGeneralPermissionBashMapItem = "deny"
)

func (r ConfigAgentGeneralPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionBashMapAsk, ConfigAgentGeneralPermissionBashMapAllow, ConfigAgentGeneralPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigAgentGeneralPermissionEdit string

const (
	ConfigAgentGeneralPermissionEditAsk   ConfigAgentGeneralPermissionEdit = "ask"
	ConfigAgentGeneralPermissionEditAllow ConfigAgentGeneralPermissionEdit = "allow"
	ConfigAgentGeneralPermissionEditDeny  ConfigAgentGeneralPermissionEdit = "deny"
)

func (r ConfigAgentGeneralPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionEditAsk, ConfigAgentGeneralPermissionEditAllow, ConfigAgentGeneralPermissionEditDeny:
		return true
	}
	return false
}

type ConfigAgentGeneralPermissionWebfetch string

const (
	ConfigAgentGeneralPermissionWebfetchAsk   ConfigAgentGeneralPermissionWebfetch = "ask"
	ConfigAgentGeneralPermissionWebfetchAllow ConfigAgentGeneralPermissionWebfetch = "allow"
	ConfigAgentGeneralPermissionWebfetchDeny  ConfigAgentGeneralPermissionWebfetch = "deny"
)

func (r ConfigAgentGeneralPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionWebfetchAsk, ConfigAgentGeneralPermissionWebfetchAllow, ConfigAgentGeneralPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigAgentPlan struct {
	// Description of when to use the agent
	Description string
	Disable     bool
	Mode        ConfigAgentPlanMode
	Model       string
	Permission  ConfigAgentPlanPermission
	Prompt      string
	Temperature float64
	Tools       map[string]bool
	TopP        float64
	ExtraFields map[string]interface{}
	JSON        configAgentPlanJSON
}

// configAgentPlanJSON contains the JSON metadata for the struct [ConfigAgentPlan]
type configAgentPlanJSON struct {
	Description apijson.Field
	Disable     apijson.Field
	Mode        apijson.Field
	Model       apijson.Field
	Permission  apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	Tools       apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgentPlan) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentPlanJSON) RawJSON() string {
	return r.raw
}

type ConfigAgentPlanMode string

const (
	ConfigAgentPlanModeSubagent ConfigAgentPlanMode = "subagent"
	ConfigAgentPlanModePrimary  ConfigAgentPlanMode = "primary"
	ConfigAgentPlanModeAll      ConfigAgentPlanMode = "all"
)

func (r ConfigAgentPlanMode) IsKnown() bool {
	switch r {
	case ConfigAgentPlanModeSubagent, ConfigAgentPlanModePrimary, ConfigAgentPlanModeAll:
		return true
	}
	return false
}

type ConfigAgentPlanPermission struct {
	Bash     ConfigAgentPlanPermissionBashUnion
	Edit     ConfigAgentPlanPermissionEdit
	Webfetch ConfigAgentPlanPermissionWebfetch
	JSON     configAgentPlanPermissionJSON
}

// configAgentPlanPermissionJSON contains the JSON metadata for the struct
// [ConfigAgentPlanPermission]
type configAgentPlanPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigAgentPlanPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configAgentPlanPermissionJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ConfigAgentPlanPermissionBashString] or
// [ConfigAgentPlanPermissionBashMap].
type ConfigAgentPlanPermissionBashUnion interface {
	implementsConfigAgentPlanPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigAgentPlanPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigAgentPlanPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigAgentPlanPermissionBashMap{}),
		},
	)
}

type ConfigAgentPlanPermissionBashString string

const (
	ConfigAgentPlanPermissionBashStringAsk   ConfigAgentPlanPermissionBashString = "ask"
	ConfigAgentPlanPermissionBashStringAllow ConfigAgentPlanPermissionBashString = "allow"
	ConfigAgentPlanPermissionBashStringDeny  ConfigAgentPlanPermissionBashString = "deny"
)

func (r ConfigAgentPlanPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionBashStringAsk, ConfigAgentPlanPermissionBashStringAllow, ConfigAgentPlanPermissionBashStringDeny:
		return true
	}
	return false
}

func (r ConfigAgentPlanPermissionBashString) implementsConfigAgentPlanPermissionBashUnion() {}

type ConfigAgentPlanPermissionBashMap map[string]ConfigAgentPlanPermissionBashMapItem

func (r ConfigAgentPlanPermissionBashMap) implementsConfigAgentPlanPermissionBashUnion() {}

type ConfigAgentPlanPermissionBashMapItem string

const (
	ConfigAgentPlanPermissionBashMapAsk   ConfigAgentPlanPermissionBashMapItem = "ask"
	ConfigAgentPlanPermissionBashMapAllow ConfigAgentPlanPermissionBashMapItem = "allow"
	ConfigAgentPlanPermissionBashMapDeny  ConfigAgentPlanPermissionBashMapItem = "deny"
)

func (r ConfigAgentPlanPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionBashMapAsk, ConfigAgentPlanPermissionBashMapAllow, ConfigAgentPlanPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigAgentPlanPermissionEdit string

const (
	ConfigAgentPlanPermissionEditAsk   ConfigAgentPlanPermissionEdit = "ask"
	ConfigAgentPlanPermissionEditAllow ConfigAgentPlanPermissionEdit = "allow"
	ConfigAgentPlanPermissionEditDeny  ConfigAgentPlanPermissionEdit = "deny"
)

func (r ConfigAgentPlanPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionEditAsk, ConfigAgentPlanPermissionEditAllow, ConfigAgentPlanPermissionEditDeny:
		return true
	}
	return false
}

type ConfigAgentPlanPermissionWebfetch string

const (
	ConfigAgentPlanPermissionWebfetchAsk   ConfigAgentPlanPermissionWebfetch = "ask"
	ConfigAgentPlanPermissionWebfetchAllow ConfigAgentPlanPermissionWebfetch = "allow"
	ConfigAgentPlanPermissionWebfetchDeny  ConfigAgentPlanPermissionWebfetch = "deny"
)

func (r ConfigAgentPlanPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionWebfetchAsk, ConfigAgentPlanPermissionWebfetchAllow, ConfigAgentPlanPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigCommand struct {
	Template    string
	Agent       string
	Description string
	Model       string
	JSON        configCommandJSON
}

// configCommandJSON contains the JSON metadata for the struct [ConfigCommand]
type configCommandJSON struct {
	Template    apijson.Field
	Agent       apijson.Field
	Description apijson.Field
	Model       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigCommand) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configCommandJSON) RawJSON() string {
	return r.raw
}

type ConfigExperimental struct {
	Hook ConfigExperimentalHook
	JSON configExperimentalJSON
}

// configExperimentalJSON contains the JSON metadata for the struct
// [ConfigExperimental]
type configExperimentalJSON struct {
	Hook        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigExperimental) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configExperimentalJSON) RawJSON() string {
	return r.raw
}

type ConfigExperimentalHook struct {
	FileEdited       map[string][]ConfigExperimentalHookFileEdited
	SessionCompleted []ConfigExperimentalHookSessionCompleted
	JSON             configExperimentalHookJSON
}

// configExperimentalHookJSON contains the JSON metadata for the struct
// [ConfigExperimentalHook]
type configExperimentalHookJSON struct {
	FileEdited       apijson.Field
	SessionCompleted apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *ConfigExperimentalHook) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configExperimentalHookJSON) RawJSON() string {
	return r.raw
}

type ConfigExperimentalHookFileEdited struct {
	Command     []string
	Environment map[string]string
	JSON        configExperimentalHookFileEditedJSON
}

// configExperimentalHookFileEditedJSON contains the JSON metadata for the struct
// [ConfigExperimentalHookFileEdited]
type configExperimentalHookFileEditedJSON struct {
	Command     apijson.Field
	Environment apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigExperimentalHookFileEdited) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configExperimentalHookFileEditedJSON) RawJSON() string {
	return r.raw
}

type ConfigExperimentalHookSessionCompleted struct {
	Command     []string
	Environment map[string]string
	JSON        configExperimentalHookSessionCompletedJSON
}

// configExperimentalHookSessionCompletedJSON contains the JSON metadata for the
// struct [ConfigExperimentalHookSessionCompleted]
type configExperimentalHookSessionCompletedJSON struct {
	Command     apijson.Field
	Environment apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigExperimentalHookSessionCompleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configExperimentalHookSessionCompletedJSON) RawJSON() string {
	return r.raw
}

type ConfigFormatter struct {
	Command     []string
	Disabled    bool
	Environment map[string]string
	Extensions  []string
	JSON        configFormatterJSON
}

// configFormatterJSON contains the JSON metadata for the struct [ConfigFormatter]
type configFormatterJSON struct {
	Command     apijson.Field
	Disabled    apijson.Field
	Environment apijson.Field
	Extensions  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigFormatter) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configFormatterJSON) RawJSON() string {
	return r.raw
}

// @deprecated Always uses stretch layout.
type ConfigLayout string

const (
	ConfigLayoutAuto    ConfigLayout = "auto"
	ConfigLayoutStretch ConfigLayout = "stretch"
)

func (r ConfigLayout) IsKnown() bool {
	switch r {
	case ConfigLayoutAuto, ConfigLayoutStretch:
		return true
	}
	return false
}

type ConfigLsp struct {
	// This field can have the runtime type of [[]string].
	Command  interface{}
	Disabled bool
	// This field can have the runtime type of [map[string]string].
	Env interface{}
	// This field can have the runtime type of [[]string].
	Extensions interface{}
	// This field can have the runtime type of [map[string]interface{}].
	Initialization interface{}
	JSON           configLspJSON
	union          ConfigLspUnion
}

// configLspJSON contains the JSON metadata for the struct [ConfigLsp]
type configLspJSON struct {
	Command        apijson.Field
	Disabled       apijson.Field
	Env            apijson.Field
	Extensions     apijson.Field
	Initialization apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r configLspJSON) RawJSON() string {
	return r.raw
}

func (r *ConfigLsp) UnmarshalJSON(data []byte) (err error) {
	*r = ConfigLsp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ConfigLspUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [ConfigLspDisabled], [ConfigLspObject].
func (r ConfigLsp) AsUnion() ConfigLspUnion {
	return r.union
}

// Union satisfied by [ConfigLspDisabled] or [ConfigLspObject].
type ConfigLspUnion interface {
	implementsConfigLsp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigLspUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigLspDisabled{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigLspObject{}),
		},
	)
}

type ConfigLspDisabled struct {
	Disabled ConfigLspDisabledDisabled
	JSON     configLspDisabledJSON
}

// configLspDisabledJSON contains the JSON metadata for the struct
// [ConfigLspDisabled]
type configLspDisabledJSON struct {
	Disabled    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigLspDisabled) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configLspDisabledJSON) RawJSON() string {
	return r.raw
}

func (r ConfigLspDisabled) implementsConfigLsp() {}

type ConfigLspDisabledDisabled bool

const (
	ConfigLspDisabledDisabledTrue ConfigLspDisabledDisabled = true
)

func (r ConfigLspDisabledDisabled) IsKnown() bool {
	switch r {
	case ConfigLspDisabledDisabledTrue:
		return true
	}
	return false
}

type ConfigLspObject struct {
	Command        []string
	Disabled       bool
	Env            map[string]string
	Extensions     []string
	Initialization map[string]interface{}
	JSON           configLspObjectJSON
}

// configLspObjectJSON contains the JSON metadata for the struct [ConfigLspObject]
type configLspObjectJSON struct {
	Command        apijson.Field
	Disabled       apijson.Field
	Env            apijson.Field
	Extensions     apijson.Field
	Initialization apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *ConfigLspObject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configLspObjectJSON) RawJSON() string {
	return r.raw
}

func (r ConfigLspObject) implementsConfigLsp() {}

type ConfigMcp struct {
	// Type of MCP server connection
	Type ConfigMcpType
	// This field can have the runtime type of [[]string].
	Command interface{}
	// Enable or disable the MCP server on startup
	Enabled bool
	// This field can have the runtime type of [map[string]string].
	Environment interface{}
	// This field can have the runtime type of [map[string]string].
	Headers interface{}
	// URL of the remote MCP server
	URL   string
	JSON  configMcpJSON
	union ConfigMcpUnion
}

// configMcpJSON contains the JSON metadata for the struct [ConfigMcp]
type configMcpJSON struct {
	Type        apijson.Field
	Command     apijson.Field
	Enabled     apijson.Field
	Environment apijson.Field
	Headers     apijson.Field
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r configMcpJSON) RawJSON() string {
	return r.raw
}

func (r *ConfigMcp) UnmarshalJSON(data []byte) (err error) {
	*r = ConfigMcp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ConfigMcpUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [McpLocalConfig], [McpRemoteConfig].
func (r ConfigMcp) AsUnion() ConfigMcpUnion {
	return r.union
}

// Union satisfied by [McpLocalConfig] or [McpRemoteConfig].
type ConfigMcpUnion interface {
	implementsConfigMcp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigMcpUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(McpLocalConfig{}),
			DiscriminatorValue: "local",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(McpRemoteConfig{}),
			DiscriminatorValue: "remote",
		},
	)
}

// Type of MCP server connection
type ConfigMcpType string

const (
	ConfigMcpTypeLocal  ConfigMcpType = "local"
	ConfigMcpTypeRemote ConfigMcpType = "remote"
)

func (r ConfigMcpType) IsKnown() bool {
	switch r {
	case ConfigMcpTypeLocal, ConfigMcpTypeRemote:
		return true
	}
	return false
}

// @deprecated Use `agent` field instead.
type ConfigMode struct {
	Build       ConfigModeBuild
	Plan        ConfigModePlan
	ExtraFields map[string]ConfigMode
	JSON        configModeJSON
}

// configModeJSON contains the JSON metadata for the struct [ConfigMode]
type configModeJSON struct {
	Build       apijson.Field
	Plan        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigMode) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configModeJSON) RawJSON() string {
	return r.raw
}

type ConfigModeBuild struct {
	// Description of when to use the agent
	Description string
	Disable     bool
	Mode        ConfigModeBuildMode
	Model       string
	Permission  ConfigModeBuildPermission
	Prompt      string
	Temperature float64
	Tools       map[string]bool
	TopP        float64
	ExtraFields map[string]interface{}
	JSON        configModeBuildJSON
}

// configModeBuildJSON contains the JSON metadata for the struct [ConfigModeBuild]
type configModeBuildJSON struct {
	Description apijson.Field
	Disable     apijson.Field
	Mode        apijson.Field
	Model       apijson.Field
	Permission  apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	Tools       apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigModeBuild) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configModeBuildJSON) RawJSON() string {
	return r.raw
}

type ConfigModeBuildMode string

const (
	ConfigModeBuildModeSubagent ConfigModeBuildMode = "subagent"
	ConfigModeBuildModePrimary  ConfigModeBuildMode = "primary"
	ConfigModeBuildModeAll      ConfigModeBuildMode = "all"
)

func (r ConfigModeBuildMode) IsKnown() bool {
	switch r {
	case ConfigModeBuildModeSubagent, ConfigModeBuildModePrimary, ConfigModeBuildModeAll:
		return true
	}
	return false
}

type ConfigModeBuildPermission struct {
	Bash     ConfigModeBuildPermissionBashUnion
	Edit     ConfigModeBuildPermissionEdit
	Webfetch ConfigModeBuildPermissionWebfetch
	JSON     configModeBuildPermissionJSON
}

// configModeBuildPermissionJSON contains the JSON metadata for the struct
// [ConfigModeBuildPermission]
type configModeBuildPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigModeBuildPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configModeBuildPermissionJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ConfigModeBuildPermissionBashString] or
// [ConfigModeBuildPermissionBashMap].
type ConfigModeBuildPermissionBashUnion interface {
	implementsConfigModeBuildPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigModeBuildPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigModeBuildPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigModeBuildPermissionBashMap{}),
		},
	)
}

type ConfigModeBuildPermissionBashString string

const (
	ConfigModeBuildPermissionBashStringAsk   ConfigModeBuildPermissionBashString = "ask"
	ConfigModeBuildPermissionBashStringAllow ConfigModeBuildPermissionBashString = "allow"
	ConfigModeBuildPermissionBashStringDeny  ConfigModeBuildPermissionBashString = "deny"
)

func (r ConfigModeBuildPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionBashStringAsk, ConfigModeBuildPermissionBashStringAllow, ConfigModeBuildPermissionBashStringDeny:
		return true
	}
	return false
}

func (r ConfigModeBuildPermissionBashString) implementsConfigModeBuildPermissionBashUnion() {}

type ConfigModeBuildPermissionBashMap map[string]ConfigModeBuildPermissionBashMapItem

func (r ConfigModeBuildPermissionBashMap) implementsConfigModeBuildPermissionBashUnion() {}

type ConfigModeBuildPermissionBashMapItem string

const (
	ConfigModeBuildPermissionBashMapAsk   ConfigModeBuildPermissionBashMapItem = "ask"
	ConfigModeBuildPermissionBashMapAllow ConfigModeBuildPermissionBashMapItem = "allow"
	ConfigModeBuildPermissionBashMapDeny  ConfigModeBuildPermissionBashMapItem = "deny"
)

func (r ConfigModeBuildPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionBashMapAsk, ConfigModeBuildPermissionBashMapAllow, ConfigModeBuildPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigModeBuildPermissionEdit string

const (
	ConfigModeBuildPermissionEditAsk   ConfigModeBuildPermissionEdit = "ask"
	ConfigModeBuildPermissionEditAllow ConfigModeBuildPermissionEdit = "allow"
	ConfigModeBuildPermissionEditDeny  ConfigModeBuildPermissionEdit = "deny"
)

func (r ConfigModeBuildPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionEditAsk, ConfigModeBuildPermissionEditAllow, ConfigModeBuildPermissionEditDeny:
		return true
	}
	return false
}

type ConfigModeBuildPermissionWebfetch string

const (
	ConfigModeBuildPermissionWebfetchAsk   ConfigModeBuildPermissionWebfetch = "ask"
	ConfigModeBuildPermissionWebfetchAllow ConfigModeBuildPermissionWebfetch = "allow"
	ConfigModeBuildPermissionWebfetchDeny  ConfigModeBuildPermissionWebfetch = "deny"
)

func (r ConfigModeBuildPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionWebfetchAsk, ConfigModeBuildPermissionWebfetchAllow, ConfigModeBuildPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigModePlan struct {
	// Description of when to use the agent
	Description string
	Disable     bool
	Mode        ConfigModePlanMode
	Model       string
	Permission  ConfigModePlanPermission
	Prompt      string
	Temperature float64
	Tools       map[string]bool
	TopP        float64
	ExtraFields map[string]interface{}
	JSON        configModePlanJSON
}

// configModePlanJSON contains the JSON metadata for the struct [ConfigModePlan]
type configModePlanJSON struct {
	Description apijson.Field
	Disable     apijson.Field
	Mode        apijson.Field
	Model       apijson.Field
	Permission  apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	Tools       apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigModePlan) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configModePlanJSON) RawJSON() string {
	return r.raw
}

type ConfigModePlanMode string

const (
	ConfigModePlanModeSubagent ConfigModePlanMode = "subagent"
	ConfigModePlanModePrimary  ConfigModePlanMode = "primary"
	ConfigModePlanModeAll      ConfigModePlanMode = "all"
)

func (r ConfigModePlanMode) IsKnown() bool {
	switch r {
	case ConfigModePlanModeSubagent, ConfigModePlanModePrimary, ConfigModePlanModeAll:
		return true
	}
	return false
}

type ConfigModePlanPermission struct {
	Bash     ConfigModePlanPermissionBashUnion
	Edit     ConfigModePlanPermissionEdit
	Webfetch ConfigModePlanPermissionWebfetch
	JSON     configModePlanPermissionJSON
}

// configModePlanPermissionJSON contains the JSON metadata for the struct
// [ConfigModePlanPermission]
type configModePlanPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigModePlanPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configModePlanPermissionJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ConfigModePlanPermissionBashString] or
// [ConfigModePlanPermissionBashMap].
type ConfigModePlanPermissionBashUnion interface {
	implementsConfigModePlanPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigModePlanPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigModePlanPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigModePlanPermissionBashMap{}),
		},
	)
}

type ConfigModePlanPermissionBashString string

const (
	ConfigModePlanPermissionBashStringAsk   ConfigModePlanPermissionBashString = "ask"
	ConfigModePlanPermissionBashStringAllow ConfigModePlanPermissionBashString = "allow"
	ConfigModePlanPermissionBashStringDeny  ConfigModePlanPermissionBashString = "deny"
)

func (r ConfigModePlanPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionBashStringAsk, ConfigModePlanPermissionBashStringAllow, ConfigModePlanPermissionBashStringDeny:
		return true
	}
	return false
}

func (r ConfigModePlanPermissionBashString) implementsConfigModePlanPermissionBashUnion() {}

type ConfigModePlanPermissionBashMap map[string]ConfigModePlanPermissionBashMapItem

func (r ConfigModePlanPermissionBashMap) implementsConfigModePlanPermissionBashUnion() {}

type ConfigModePlanPermissionBashMapItem string

const (
	ConfigModePlanPermissionBashMapAsk   ConfigModePlanPermissionBashMapItem = "ask"
	ConfigModePlanPermissionBashMapAllow ConfigModePlanPermissionBashMapItem = "allow"
	ConfigModePlanPermissionBashMapDeny  ConfigModePlanPermissionBashMapItem = "deny"
)

func (r ConfigModePlanPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionBashMapAsk, ConfigModePlanPermissionBashMapAllow, ConfigModePlanPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigModePlanPermissionEdit string

const (
	ConfigModePlanPermissionEditAsk   ConfigModePlanPermissionEdit = "ask"
	ConfigModePlanPermissionEditAllow ConfigModePlanPermissionEdit = "allow"
	ConfigModePlanPermissionEditDeny  ConfigModePlanPermissionEdit = "deny"
)

func (r ConfigModePlanPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionEditAsk, ConfigModePlanPermissionEditAllow, ConfigModePlanPermissionEditDeny:
		return true
	}
	return false
}

type ConfigModePlanPermissionWebfetch string

const (
	ConfigModePlanPermissionWebfetchAsk   ConfigModePlanPermissionWebfetch = "ask"
	ConfigModePlanPermissionWebfetchAllow ConfigModePlanPermissionWebfetch = "allow"
	ConfigModePlanPermissionWebfetchDeny  ConfigModePlanPermissionWebfetch = "deny"
)

func (r ConfigModePlanPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionWebfetchAsk, ConfigModePlanPermissionWebfetchAllow, ConfigModePlanPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigPermission struct {
	Bash     ConfigPermissionBashUnion
	Edit     ConfigPermissionEdit
	Webfetch ConfigPermissionWebfetch
	JSON     configPermissionJSON
}

// configPermissionJSON contains the JSON metadata for the struct
// [ConfigPermission]
type configPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configPermissionJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ConfigPermissionBashString] or [ConfigPermissionBashMap].
type ConfigPermissionBashUnion interface {
	implementsConfigPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigPermissionBashMap{}),
		},
	)
}

type ConfigPermissionBashString string

const (
	ConfigPermissionBashStringAsk   ConfigPermissionBashString = "ask"
	ConfigPermissionBashStringAllow ConfigPermissionBashString = "allow"
	ConfigPermissionBashStringDeny  ConfigPermissionBashString = "deny"
)

func (r ConfigPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigPermissionBashStringAsk, ConfigPermissionBashStringAllow, ConfigPermissionBashStringDeny:
		return true
	}
	return false
}

func (r ConfigPermissionBashString) implementsConfigPermissionBashUnion() {}

type ConfigPermissionBashMap map[string]ConfigPermissionBashMapItem

func (r ConfigPermissionBashMap) implementsConfigPermissionBashUnion() {}

type ConfigPermissionBashMapItem string

const (
	ConfigPermissionBashMapAsk   ConfigPermissionBashMapItem = "ask"
	ConfigPermissionBashMapAllow ConfigPermissionBashMapItem = "allow"
	ConfigPermissionBashMapDeny  ConfigPermissionBashMapItem = "deny"
)

func (r ConfigPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigPermissionBashMapAsk, ConfigPermissionBashMapAllow, ConfigPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigPermissionEdit string

const (
	ConfigPermissionEditAsk   ConfigPermissionEdit = "ask"
	ConfigPermissionEditAllow ConfigPermissionEdit = "allow"
	ConfigPermissionEditDeny  ConfigPermissionEdit = "deny"
)

func (r ConfigPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigPermissionEditAsk, ConfigPermissionEditAllow, ConfigPermissionEditDeny:
		return true
	}
	return false
}

type ConfigPermissionWebfetch string

const (
	ConfigPermissionWebfetchAsk   ConfigPermissionWebfetch = "ask"
	ConfigPermissionWebfetchAllow ConfigPermissionWebfetch = "allow"
	ConfigPermissionWebfetchDeny  ConfigPermissionWebfetch = "deny"
)

func (r ConfigPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigPermissionWebfetchAsk, ConfigPermissionWebfetchAllow, ConfigPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigProvider struct {
	ID      string
	API     string
	Env     []string
	Models  map[string]ConfigProviderModel
	Name    string
	Npm     string
	Options ConfigProviderOptions
	JSON    configProviderJSON
}

// configProviderJSON contains the JSON metadata for the struct [ConfigProvider]
type configProviderJSON struct {
	ID          apijson.Field
	API         apijson.Field
	Env         apijson.Field
	Models      apijson.Field
	Name        apijson.Field
	Npm         apijson.Field
	Options     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigProvider) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configProviderJSON) RawJSON() string {
	return r.raw
}

type ConfigProviderModel struct {
	ID          string
	Attachment  bool
	Cost        ConfigProviderModelsCost
	Limit       ConfigProviderModelsLimit
	Name        string
	Options     map[string]interface{}
	Reasoning   bool
	ReleaseDate string
	Temperature bool
	ToolCall    bool
	JSON        configProviderModelJSON
}

// configProviderModelJSON contains the JSON metadata for the struct
// [ConfigProviderModel]
type configProviderModelJSON struct {
	ID          apijson.Field
	Attachment  apijson.Field
	Cost        apijson.Field
	Limit       apijson.Field
	Name        apijson.Field
	Options     apijson.Field
	Reasoning   apijson.Field
	ReleaseDate apijson.Field
	Temperature apijson.Field
	ToolCall    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigProviderModel) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configProviderModelJSON) RawJSON() string {
	return r.raw
}

type ConfigProviderModelsCost struct {
	Input      float64
	Output     float64
	CacheRead  float64
	CacheWrite float64
	JSON       configProviderModelsCostJSON
}

// configProviderModelsCostJSON contains the JSON metadata for the struct
// [ConfigProviderModelsCost]
type configProviderModelsCostJSON struct {
	Input       apijson.Field
	Output      apijson.Field
	CacheRead   apijson.Field
	CacheWrite  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigProviderModelsCost) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configProviderModelsCostJSON) RawJSON() string {
	return r.raw
}

type ConfigProviderModelsLimit struct {
	Context float64
	Output  float64
	JSON    configProviderModelsLimitJSON
}

// configProviderModelsLimitJSON contains the JSON metadata for the struct
// [ConfigProviderModelsLimit]
type configProviderModelsLimitJSON struct {
	Context     apijson.Field
	Output      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigProviderModelsLimit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configProviderModelsLimitJSON) RawJSON() string {
	return r.raw
}

type ConfigProviderOptions struct {
	APIKey  string
	BaseURL string
	// Timeout in milliseconds for requests to this provider. Default is 300000 (5
	// minutes). Set to false to disable timeout.
	Timeout     ConfigProviderOptionsTimeoutUnion
	ExtraFields map[string]interface{}
	JSON        configProviderOptionsJSON
}

// configProviderOptionsJSON contains the JSON metadata for the struct
// [ConfigProviderOptions]
type configProviderOptionsJSON struct {
	APIKey      apijson.Field
	BaseURL     apijson.Field
	Timeout     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigProviderOptions) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configProviderOptionsJSON) RawJSON() string {
	return r.raw
}

// Timeout in milliseconds for requests to this provider. Default is 300000 (5
// minutes). Set to false to disable timeout.
//
// Union satisfied by [shared.UnionInt] or [shared.UnionBool].
type ConfigProviderOptionsTimeoutUnion interface {
	ImplementsConfigProviderOptionsTimeoutUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigProviderOptionsTimeoutUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.True,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.False,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
	)
}

// Control sharing behavior:'manual' allows manual sharing via commands, 'auto'
// enables automatic sharing, 'disabled' disables all sharing
type ConfigShare string

const (
	ConfigShareManual   ConfigShare = "manual"
	ConfigShareAuto     ConfigShare = "auto"
	ConfigShareDisabled ConfigShare = "disabled"
)

func (r ConfigShare) IsKnown() bool {
	switch r {
	case ConfigShareManual, ConfigShareAuto, ConfigShareDisabled:
		return true
	}
	return false
}

// TUI specific settings
type ConfigTui struct {
	// TUI scroll speed
	ScrollSpeed float64
	JSON        configTuiJSON
}

// configTuiJSON contains the JSON metadata for the struct [ConfigTui]
type configTuiJSON struct {
	ScrollSpeed apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ConfigTui) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r configTuiJSON) RawJSON() string {
	return r.raw
}

type KeybindsConfig struct {
	// Next agent
	AgentCycle string
	// Previous agent
	AgentCycleReverse string
	// List agents
	AgentList string
	// Exit the application
	AppExit string
	// Show help dialog
	AppHelp string
	// Open external editor
	EditorOpen string
	// @deprecated Close file
	FileClose string
	// @deprecated Split/unified diff
	FileDiffToggle string
	// @deprecated Currently not available. List files
	FileList string
	// @deprecated Search file
	FileSearch string
	// Clear input field
	InputClear string
	// Insert newline in input
	InputNewline string
	// Paste from clipboard
	InputPaste string
	// Submit input
	InputSubmit string
	// Leader key for keybind combinations
	Leader string
	// Copy message
	MessagesCopy string
	// Navigate to first message
	MessagesFirst string
	// Scroll messages down by half page
	MessagesHalfPageDown string
	// Scroll messages up by half page
	MessagesHalfPageUp string
	// Navigate to last message
	MessagesLast string
	// @deprecated Toggle layout
	MessagesLayoutToggle string
	// @deprecated Navigate to next message
	MessagesNext string
	// Scroll messages down by one page
	MessagesPageDown string
	// Scroll messages up by one page
	MessagesPageUp string
	// @deprecated Navigate to previous message
	MessagesPrevious string
	// Redo message
	MessagesRedo string
	// @deprecated use messages_undo. Revert message
	MessagesRevert string
	// Undo message
	MessagesUndo string
	// Next recent model
	ModelCycleRecent string
	// Previous recent model
	ModelCycleRecentReverse string
	// List available models
	ModelList string
	// Create/update AGENTS.md
	ProjectInit string
	// Cycle to next child session
	SessionChildCycle string
	// Cycle to previous child session
	SessionChildCycleReverse string
	// Compact the session
	SessionCompact string
	// Export session to editor
	SessionExport string
	// Interrupt current session
	SessionInterrupt string
	// List all sessions
	SessionList string
	// Create a new session
	SessionNew string
	// Share current session
	SessionShare string
	// Show session timeline
	SessionTimeline string
	// Unshare current session
	SessionUnshare string
	// @deprecated use agent_cycle. Next agent
	SwitchAgent string
	// @deprecated use agent_cycle_reverse. Previous agent
	SwitchAgentReverse string
	// @deprecated use agent_cycle. Next mode
	SwitchMode string
	// @deprecated use agent_cycle_reverse. Previous mode
	SwitchModeReverse string
	// List available themes
	ThemeList string
	// Toggle thinking blocks
	ThinkingBlocks string
	// Toggle tool details
	ToolDetails string
	JSON        keybindsConfigJSON
}

// keybindsConfigJSON contains the JSON metadata for the struct [KeybindsConfig]
type keybindsConfigJSON struct {
	AgentCycle               apijson.Field
	AgentCycleReverse        apijson.Field
	AgentList                apijson.Field
	AppExit                  apijson.Field
	AppHelp                  apijson.Field
	EditorOpen               apijson.Field
	FileClose                apijson.Field
	FileDiffToggle           apijson.Field
	FileList                 apijson.Field
	FileSearch               apijson.Field
	InputClear               apijson.Field
	InputNewline             apijson.Field
	InputPaste               apijson.Field
	InputSubmit              apijson.Field
	Leader                   apijson.Field
	MessagesCopy             apijson.Field
	MessagesFirst            apijson.Field
	MessagesHalfPageDown     apijson.Field
	MessagesHalfPageUp       apijson.Field
	MessagesLast             apijson.Field
	MessagesLayoutToggle     apijson.Field
	MessagesNext             apijson.Field
	MessagesPageDown         apijson.Field
	MessagesPageUp           apijson.Field
	MessagesPrevious         apijson.Field
	MessagesRedo             apijson.Field
	MessagesRevert           apijson.Field
	MessagesUndo             apijson.Field
	ModelCycleRecent         apijson.Field
	ModelCycleRecentReverse  apijson.Field
	ModelList                apijson.Field
	ProjectInit              apijson.Field
	SessionChildCycle        apijson.Field
	SessionChildCycleReverse apijson.Field
	SessionCompact           apijson.Field
	SessionExport            apijson.Field
	SessionInterrupt         apijson.Field
	SessionList              apijson.Field
	SessionNew               apijson.Field
	SessionShare             apijson.Field
	SessionTimeline          apijson.Field
	SessionUnshare           apijson.Field
	SwitchAgent              apijson.Field
	SwitchAgentReverse       apijson.Field
	SwitchMode               apijson.Field
	SwitchModeReverse        apijson.Field
	ThemeList                apijson.Field
	ThinkingBlocks           apijson.Field
	ToolDetails              apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *KeybindsConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r keybindsConfigJSON) RawJSON() string {
	return r.raw
}

type McpLocalConfig struct {
	// Command and arguments to run the MCP server
	Command []string
	// Type of MCP server connection
	Type McpLocalConfigType
	// Enable or disable the MCP server on startup
	Enabled bool
	// Environment variables to set when running the MCP server
	Environment map[string]string
	JSON        mcpLocalConfigJSON
}

// mcpLocalConfigJSON contains the JSON metadata for the struct [McpLocalConfig]
type mcpLocalConfigJSON struct {
	Command     apijson.Field
	Type        apijson.Field
	Enabled     apijson.Field
	Environment apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *McpLocalConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r mcpLocalConfigJSON) RawJSON() string {
	return r.raw
}

func (r McpLocalConfig) implementsConfigMcp() {}

// Type of MCP server connection
type McpLocalConfigType string

const (
	McpLocalConfigTypeLocal McpLocalConfigType = "local"
)

func (r McpLocalConfigType) IsKnown() bool {
	switch r {
	case McpLocalConfigTypeLocal:
		return true
	}
	return false
}

type McpRemoteConfig struct {
	// Type of MCP server connection
	Type McpRemoteConfigType
	// URL of the remote MCP server
	URL string
	// Enable or disable the MCP server on startup
	Enabled bool
	// Headers to send with the request
	Headers map[string]string
	JSON    mcpRemoteConfigJSON
}

// mcpRemoteConfigJSON contains the JSON metadata for the struct [McpRemoteConfig]
type mcpRemoteConfigJSON struct {
	Type        apijson.Field
	URL         apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *McpRemoteConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r mcpRemoteConfigJSON) RawJSON() string {
	return r.raw
}

func (r McpRemoteConfig) implementsConfigMcp() {}

// Type of MCP server connection
type McpRemoteConfigType string

const (
	McpRemoteConfigTypeRemote McpRemoteConfigType = "remote"
)

func (r McpRemoteConfigType) IsKnown() bool {
	switch r {
	case McpRemoteConfigTypeRemote:
		return true
	}
	return false
}

type ConfigGetParams struct {
	Directory param.Field[string]
}

// URLQuery serializes [ConfigGetParams]'s query parameters as `url.Values`.
func (r ConfigGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
