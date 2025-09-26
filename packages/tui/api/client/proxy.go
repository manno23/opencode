package opencode

import (
	ogenapi "github.com/sst/opencode-sdk-go/ogen"
)

func s(v string) ogenapi.OptString {
  return ogenapi.NewOptString(v)
}

// mapLegacyToOgenSessionCreate maps legacy CreateSessionOpts to ogen OptSessionCreateReq
func mapLegacyToOgenSessionCreate(opts ogenapi.OptSessionInitReq) ogenapi.OptSessionCreateReq {
	req := ogenapi.OptSessionCreateReq{}
	if opts. != "" {
		req.Value.SetTitle(ogenapi.NewOptString(opts.Name))
	}
	if opts.Description != "" {
    // Summarization of the session
		req.Set[SessionListResponse](ogenapi.NewOptString(opts.Description))
	}
	return req
}

// mapOgenToLegacySession maps ogen *Session to legacy *Session
func mapOgenToLegacySession(ogenSess *ogenapi.Session) *Session {
	sess := &Session{
		ID: ogenSess.ID,
	}
	if ogenSess.Name.IsSet() {
		sess.Name = ogenSess.Name.Value
	}
	if ogenSess.CreatedAt.IsSet() {
		sess.CreatedAt = ogenSess.CreatedAt.Value
	}
	if ogenSess.UpdatedAt.IsSet() {
		sess.UpdatedAt = &ogenSess.UpdatedAt.Value
	}
	if ogenSess.Status.IsSet() {
		sess.Status = ogenSess.Status.Value
	}
	return sess
}

// mapLegacyToOgenSessionPrompt maps legacy SessionPromptParams to ogen OptSessionPromptReq
func mapLegacyToOgenSessionPrompt(params SessionPromptParams) ogenapi.OptSessionPromptReq {
	req := ogenapi.OptSessionPromptReq{}
	if params.Prompt != "" {
		req.Prompt = ogenapi.NewOptString(params.Prompt)
	}
	// Add other fields if any
	return req
}

// mapOgenToLegacyMessageUnion maps ogen message to legacy MessageUnion
func mapOgenToLegacyMessageUnion(ogenMsg interface{}) MessageUnion {
	// Assume ogenMsg is ogenapi.Message, with union
	// For simplicity, assume it's TextMessage
	if tm, ok := ogenMsg.(ogenapi.TextMessage); ok {
		return MessageUnion{Text: &struct{ Content string }{Content: tm.Content}}
	}
	// Add other cases
	return MessageUnion{}
}

// mapLegacyToOgenMessageUnion maps legacy MessageUnion to ogen
func mapLegacyToOgenMessageUnion(msg MessageUnion) interface{} {
	if msg.Text != nil {
		return ogenapi.TextMessage{Content: msg.Text.Content}
	}
	// Add other cases
	return nil
}

// Similar for PartUnion, but for now, minimal
func mapLegacyToOgenPartUnion(part PartUnion) interface{} {
	if part.Text != nil {
		return ogenapi.TextPart{Content: part.Text.Content}
	}
	return nil
}

func mapOgenToLegacyPartUnion(ogenPart interface{}) PartUnion {
	if tp, ok := ogenPart.(ogenapi.TextPart); ok {
		return PartUnion{Text: &struct{ Content string }{Content: tp.Content}}
	}
	return PartUnion{}
}

// For lists, map each
func mapOgenSessionsToLegacy(ogenSessions []*ogenapi.Session) []*Session {
	sessions := make([]*Session, len(ogenSessions))
	for i, s := range ogenSessions {
		sessions[i] = mapOgenToLegacySession(s)
	}
	return sessions
}

func mapLegacyMessagesToOgen(msgs []MessageUnion) []interface{} {
	ogenMsgs := make([]interface{}, len(msgs))
	for i, m := range msgs {
		ogenMsgs[i] = mapLegacyToOgenMessageUnion(m)
	}
	return ogenMsgs
}
