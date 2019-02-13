package cloud

import "strings"

type tpLinkChildIDs struct {
	ChildIDs []string `json:"child_ids"`
}

type tpLinkContext struct {
	TPLinkChildren tpLinkChildIDs `json:"context"`
	SystemPayload  string         `json:"system"`
}

func buildContextPayload(systemCall string, ids []string) tpLinkContext {
	var c = tpLinkContext{TPLinkChildren: tpLinkChildIDs{ids}}
	c.SystemPayload = strings.TrimPrefix(systemCall, `{"system"`)
	c.SystemPayload = strings.TrimSuffix(c.SystemPayload, `}`)
	return c
}
