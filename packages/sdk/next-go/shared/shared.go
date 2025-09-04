// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"github.com/sst/opencode-sdk-go/internal/apijson"
)

type MessageAbortedError struct {
	Data interface{}
	Name MessageAbortedErrorName
	JSON messageAbortedErrorJSON
}

// messageAbortedErrorJSON contains the JSON metadata for the struct
// [MessageAbortedError]
type messageAbortedErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageAbortedError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageAbortedErrorJSON) RawJSON() string {
	return r.raw
}

func (r MessageAbortedError) ImplementsEventListResponseEventSessionErrorPropertiesError() {}

func (r MessageAbortedError) ImplementsAssistantMessageError() {}

type MessageAbortedErrorName string

const (
	MessageAbortedErrorNameMessageAbortedError MessageAbortedErrorName = "MessageAbortedError"
)

func (r MessageAbortedErrorName) IsKnown() bool {
	switch r {
	case MessageAbortedErrorNameMessageAbortedError:
		return true
	}
	return false
}

type ProviderAuthError struct {
	Data ProviderAuthErrorData
	Name ProviderAuthErrorName
	JSON providerAuthErrorJSON
}

// providerAuthErrorJSON contains the JSON metadata for the struct
// [ProviderAuthError]
type providerAuthErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProviderAuthError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r providerAuthErrorJSON) RawJSON() string {
	return r.raw
}

func (r ProviderAuthError) ImplementsEventListResponseEventSessionErrorPropertiesError() {}

func (r ProviderAuthError) ImplementsAssistantMessageError() {}

type ProviderAuthErrorData struct {
	Message    string
	ProviderID string
	JSON       providerAuthErrorDataJSON
}

// providerAuthErrorDataJSON contains the JSON metadata for the struct
// [ProviderAuthErrorData]
type providerAuthErrorDataJSON struct {
	Message     apijson.Field
	ProviderID  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProviderAuthErrorData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r providerAuthErrorDataJSON) RawJSON() string {
	return r.raw
}

type ProviderAuthErrorName string

const (
	ProviderAuthErrorNameProviderAuthError ProviderAuthErrorName = "ProviderAuthError"
)

func (r ProviderAuthErrorName) IsKnown() bool {
	switch r {
	case ProviderAuthErrorNameProviderAuthError:
		return true
	}
	return false
}

type UnknownError struct {
	Data UnknownErrorData
	Name UnknownErrorName
	JSON unknownErrorJSON
}

// unknownErrorJSON contains the JSON metadata for the struct [UnknownError]
type unknownErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *UnknownError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r unknownErrorJSON) RawJSON() string {
	return r.raw
}

func (r UnknownError) ImplementsEventListResponseEventSessionErrorPropertiesError() {}

func (r UnknownError) ImplementsAssistantMessageError() {}

type UnknownErrorData struct {
	Message string
	JSON    unknownErrorDataJSON
}

// unknownErrorDataJSON contains the JSON metadata for the struct
// [UnknownErrorData]
type unknownErrorDataJSON struct {
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *UnknownErrorData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r unknownErrorDataJSON) RawJSON() string {
	return r.raw
}

type UnknownErrorName string

const (
	UnknownErrorNameUnknownError UnknownErrorName = "UnknownError"
)

func (r UnknownErrorName) IsKnown() bool {
	switch r {
	case UnknownErrorNameUnknownError:
		return true
	}
	return false
}
