package opencode

import (
	"git.j9xym.com/opencode-api-go/internal/apierror"
	"git.j9xym.com/opencode-api-go/shared"
)

type Error = apierror.Error

// This is an alias to an internal type.
type MessageAbortedError = shared.MessageAbortedError

// This is an alias to an internal type.
type MessageAbortedErrorName = shared.MessageAbortedErrorName

// This is an alias to an internal value.
const MessageAbortedErrorNameMessageAbortedError = shared.MessageAbortedErrorNameMessageAbortedError

// This is an alias to an internal type.
type ProviderAuthError = shared.ProviderAuthError

// This is an alias to an internal type.
type ProviderAuthErrorData = shared.ProviderAuthErrorData

// This is an alias to an internal type.
type ProviderAuthErrorName = shared.ProviderAuthErrorName

// This is an alias to an internal value.
const ProviderAuthErrorNameProviderAuthError = shared.ProviderAuthErrorNameProviderAuthError

// This is an alias to an internal type.
type UnknownError = shared.UnknownError

// This is an alias to an internal type.
type UnknownErrorData = shared.UnknownErrorData

// This is an alias to an internal type.
type UnknownErrorName = shared.UnknownErrorName

// This is an alias to an internal value.
const UnknownErrorNameUnknownError = shared.UnknownErrorNameUnknownError
