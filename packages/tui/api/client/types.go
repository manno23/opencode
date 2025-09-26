package opencode

import (
	legacy "github.com/sst/opencode/sdk/go"
)

// Types to match legacy
type Session = legacy.Session

type SessionListParams = legacy.SessionListParams
type SessionListResponse struct {
	Sessions []Session `json:"sessions"`
}
type SessionCreateParams = legacy.SessionNewParams
type SessionCreateResponse struct {
	Session Session `json:"session"`
}
type SessionCloseParams = legacy.SessionAbortParams
type SessionCloseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Type aliases from legacy SDK

type Project = legacy.Project
type Path = legacy.Path
type Agent = legacy.Agent
type Provider = legacy.Provider
type Model = legacy.Model
type MessageUnion = legacy.MessageUnion
type PartUnion = legacy.PartUnion
type Permission = legacy.Permission
type Config = legacy.Config
type File = legacy.File
type Command = legacy.Command

type Tui = legacy.TuiService
type Find = legacy.FindService
type App = legacy.AppService
type Event = legacy.EventService

type UserMessage = legacy.UserMessage

// Added missing symbols for adapter
type OptSessionListReq = legacy.SessionListParams
type LegacySessionListResponse = []legacy.Session
type LegacyError = legacy.Error
