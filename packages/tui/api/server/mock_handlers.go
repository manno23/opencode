package server

import (
	"context"
	"strings"
	"time"
)

type MockServerImpl struct{}

func (m *MockServerImpl) AppGet(ctx context.Context) (*App, error) {
	return &App{
		Hostname: "localhost",
		Git:      true,
		Path: AppPath{
			Config: "/config",
			Data:   "/data",
			Root:   "/",
			Cwd:    "/cwd",
			State:  "/state",
		},
		Time: AppTime{
			Initialized: NewOptFloat64(float64(time.Now().Unix())),
		},
	}, nil
}

func (m *MockServerImpl) AppInit(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) AppLog(ctx context.Context, req OptAppLogReq) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) AuthSet(ctx context.Context, req OptAuth, params AuthSetParams) (AuthSetRes, error) {
	return &Error{
		Data: ErrorData{},
	}, nil
}

func (m *MockServerImpl) CommandList(ctx context.Context) ([]Command, error) {
	return []Command{}, nil
}

func (m *MockServerImpl) ConfigProviders(ctx context.Context) (*ConfigProvidersOK, error) {
	return &ConfigProvidersOK{
		Providers: []Provider{},
	}, nil
}

func (m *MockServerImpl) EventSubscribe(ctx context.Context) (EventSubscribeOK, error) {
	return EventSubscribeOK{Data: strings.NewReader("{}")}, nil
}

func (m *MockServerImpl) FileRead(ctx context.Context, params FileReadParams) (*FileReadOK, error) {
	return &FileReadOK{
		Content: "test content",
	}, nil
}

func (m *MockServerImpl) FileStatus(ctx context.Context) ([]File, error) {
	return []File{}, nil
}

func (m *MockServerImpl) FindFiles(ctx context.Context, params FindFilesParams) ([]string, error) {
	return []string{}, nil
}

func (m *MockServerImpl) FindSymbols(ctx context.Context, params FindSymbolsParams) ([]Symbol, error) {
	return []Symbol{}, nil
}

func (m *MockServerImpl) FindText(ctx context.Context, params FindTextParams) ([]FindTextOKItem, error) {
	return []FindTextOKItem{}, nil
}

func (m *MockServerImpl) PostSessionByIdPermissionsByPermissionID(ctx context.Context, req OptPostSessionByIdPermissionsByPermissionIDReq, params PostSessionByIdPermissionsByPermissionIDParams) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) SessionAbort(ctx context.Context, params SessionAbortParams) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) SessionChat(ctx context.Context, req OptSessionChatReq, params SessionChatParams) (*SessionChatOK, error) {
	return &SessionChatOK{
		Info: AssistantMessage{
			ID:        "msg-1",
			SessionID: "test-session",
			Role:      "assistant",
			Time: AssistantMessageTime{
				Created: float64(time.Now().Unix()),
			},
			System:     []string{},
			ModelID:    "model",
			ProviderID: "prov",
			Mode:       "chat",
			Path: AssistantMessagePath{
				Cwd:  "",
				Root: "",
			},
			Cost: 0,
			Tokens: AssistantMessageTokens{
				Input:     0,
				Output:    0,
				Reasoning: 0,
				Cache: AssistantMessageTokensCache{
					Read:  0,
					Write: 0,
				},
			},
		},
		Parts: []Part{},
	}, nil
}

func (m *MockServerImpl) SessionCommand(ctx context.Context, req OptSessionCommandReq, params SessionCommandParams) (*SessionCommandOK, error) {
	return &SessionCommandOK{
		Info: AssistantMessage{
			ID:        "msg-1",
			SessionID: "test-session",
			Role:      "assistant",
			Time: AssistantMessageTime{
				Created: float64(time.Now().Unix()),
			},
			System:     []string{},
			ModelID:    "model",
			ProviderID: "prov",
			Mode:       "chat",
			Path: AssistantMessagePath{
				Cwd:  "",
				Root: "",
			},
			Cost: 0,
			Tokens: AssistantMessageTokens{
				Input:     0,
				Output:    0,
				Reasoning: 0,
				Cache: AssistantMessageTokensCache{
					Read:  0,
					Write: 0,
				},
			},
		},
		Parts: []Part{},
	}, nil
}

func (m *MockServerImpl) SessionCreate(ctx context.Context, req OptSessionCreateReq) (SessionCreateRes, error) {
	return &Session{
		ID:      "new-session",
		Title:   "New Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) SessionDelete(ctx context.Context, params SessionDeleteParams) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) SessionGet(ctx context.Context, params SessionGetParams) (*Session, error) {
	return &Session{
		ID:      "test-session",
		Title:   "Test Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) SessionInit(ctx context.Context, req OptSessionInitReq, params SessionInitParams) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) SessionMessage(ctx context.Context, params SessionMessageParams) (*SessionMessageOK, error) {
	return &SessionMessageOK{
		Info: Message{
			Type: AssistantMessageMessage,
			AssistantMessage: AssistantMessage{
				ID:        "msg-1",
				SessionID: "test-session",
				Role:      "assistant",
				Time: AssistantMessageTime{
					Created: float64(time.Now().Unix()),
				},
				System:     []string{},
				ModelID:    "model",
				ProviderID: "prov",
				Mode:       "chat",
				Path: AssistantMessagePath{
					Cwd:  "",
					Root: "",
				},
				Cost: 0,
				Tokens: AssistantMessageTokens{
					Input:     0,
					Output:    0,
					Reasoning: 0,
					Cache: AssistantMessageTokensCache{
						Read:  0,
						Write: 0,
					},
				},
			},
		},
		Parts: []Part{},
	}, nil
}

func (m *MockServerImpl) SessionList(ctx context.Context) ([]Session, error) {
	return []Session{
		{
			ID:      "session-1",
			Title:   "Session 1",
			Version: "1.0",
			Time: SessionTime{
				Created: float64(time.Now().Unix()),
			},
		},
	}, nil
}

func (m *MockServerImpl) SessionChildren(ctx context.Context, params SessionChildrenParams) ([]Session, error) {
	return []Session{}, nil
}

func (m *MockServerImpl) SessionMessages(ctx context.Context, params SessionMessagesParams) ([]SessionMessagesOKItem, error) {
	return []SessionMessagesOKItem{}, nil
}

func (m *MockServerImpl) SessionRevert(ctx context.Context, req OptSessionRevertReq, params SessionRevertParams) (*Session, error) {
	return &Session{
		ID:      "test-session",
		Title:   "Test Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) SessionShare(ctx context.Context, params SessionShareParams) (*Session, error) {
	return &Session{
		ID:      "test-session",
		Title:   "Test Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) SessionShell(ctx context.Context, req OptSessionShellReq, params SessionShellParams) (*AssistantMessage, error) {
	return &AssistantMessage{
		ID:        "msg-1",
		SessionID: "test-session",
		Role:      "assistant",
		Time: AssistantMessageTime{
			Created: float64(time.Now().Unix()),
		},
		System:     []string{},
		ModelID:    "model",
		ProviderID: "prov",
		Mode:       "chat",
		Path: AssistantMessagePath{
			Cwd:  "",
			Root: "",
		},
		Cost: 0,
		Tokens: AssistantMessageTokens{
			Input:     0,
			Output:    0,
			Reasoning: 0,
			Cache: AssistantMessageTokensCache{
				Read:  0,
				Write: 0,
			},
		},
	}, nil
}

func (m *MockServerImpl) SessionSummarize(ctx context.Context, req OptSessionSummarizeReq, params SessionSummarizeParams) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) SessionUnrevert(ctx context.Context, params SessionUnrevertParams) (*Session, error) {
	return &Session{
		ID:      "test-session",
		Title:   "Test Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) SessionUnshare(ctx context.Context, params SessionUnshareParams) (*Session, error) {
	return &Session{
		ID:      "test-session",
		Title:   "Test Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) SessionUpdate(ctx context.Context, req OptSessionUpdateReq, params SessionUpdateParams) (*Session, error) {
	return &Session{
		ID:      "test-session",
		Title:   "Test Session",
		Version: "1.0",
		Time: SessionTime{
			Created: float64(time.Now().Unix()),
		},
	}, nil
}

func (m *MockServerImpl) TuiAppendPrompt(ctx context.Context, req OptTuiAppendPromptReq) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiClearPrompt(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiExecuteCommand(ctx context.Context, req OptTuiExecuteCommandReq) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiOpenHelp(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiOpenModels(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiOpenSessions(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiOpenThemes(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiShowToast(ctx context.Context, req OptTuiShowToastReq) (bool, error) {
	return true, nil
}

func (m *MockServerImpl) TuiSubmitPrompt(ctx context.Context) (bool, error) {
	return true, nil
}
