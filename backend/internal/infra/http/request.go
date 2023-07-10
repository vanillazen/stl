package http

import (
	"net/http"
)

type ContextKey string

const (
	ReqCtxKey  = "req"
	UserCtxKey = "user"
	ListCtxKey = "list"
)

type AssetRequest struct {
	Type        string `json:"type"`
	Action      string `json:"action"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

const (
	UserIDCtxKey   = "resource"
	ResourceCtxKey = "resource"
	AssetReqCtxKey = "assetreq"
)

func (h *APIHandler) userID(r *http.Request) (userID string, ok bool) {
	value := r.Context().Value(UserIDCtxKey)
	if value == nil {
		return userID, false
	}

	userID, ok = value.(string)
	if !ok {
		return userID, false
	}

	return userID, true
}

func (h *APIHandler) resource(r *http.Request) (ri ResourceInfo, ok bool) {
	value := r.Context().Value(ResourceCtxKey)
	if value == nil {
		return ri, false
	}

	ri, ok = value.(ResourceInfo)
	if !ok {
		return ri, false
	}

	return ri, true
}

func (h *APIHandler) assetReq(r *http.Request) (req AssetRequest, ok bool) {
	value := r.Context().Value(AssetReqCtxKey)
	if value == nil {
		return req, false
	}

	req, ok = value.(AssetRequest)
	if !ok {
		return req, false
	}

	return req, true
}
