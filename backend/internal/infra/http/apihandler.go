package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vanillazen/stl/backend/internal/domain/service"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type (
	APIHTTPHandler interface {
		sys.Core
		Service() service.ListService
	}

	APIHandler struct {
		*sys.SimpleCore
		svc    service.ListService
		apiDoc string
	}
)

func NewAPIHandler(svc service.ListService, apiDoc string, opts ...sys.Option) *APIHandler {
	return &APIHandler{
		SimpleCore: sys.NewCore("list-handler", opts...),
		svc:        svc,
		apiDoc:     apiDoc,
	}
}

func (h *APIHandler) handleV1(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	userIDIndex := -1
	for i, segment := range pathSegments {
		ok := uuid.Validate(segment)
		if ok {
			userIDIndex = i
			break
		}
	}

	if userIDIndex == -1 {
		http.NotFound(w, r)
		return
	}

	userIDSegment := pathSegments[userIDIndex]
	ctx := context.WithValue(r.Context(), UserIDCtxKey, userIDSegment)
	r = r.WithContext(ctx)

	resIDSegment := pathSegments[len(pathSegments)-1]
	resID, err := uuid.Parse(resIDSegment)
	if err != nil {
		h.Log().Debug("Not a resource URL:", r.URL.Path)
	} else {
		ctx = context.WithValue(r.Context(), ResIDCtxKey, resID)
		r = r.WithContext(ctx)
	}

	resource := strings.ToLower(pathSegments[userIDIndex+1])

	switch resource {
	case "list":
		h.handleList(w, r)
	default:
		h.handleError(w, InvalidResourceErr)
	}
}

func (h *APIHandler) handleList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, ok := h.resourceID(r)

		switch ok {
		case true:
			h.GetList(w, r)
			h.Log().Error("not implemented yet")

		default:
			//h.GetLists(w, r)
			h.Log().Error("not implemented yet")
		}

	case http.MethodPost:
		//h.CreateListList(w, r)
		h.Log().Error("not implemented yet")

	case http.MethodPut:
		//h.UpdateList(w, r)
		h.Log().Error("not implemented yet")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		h.handleError(w, MethodNotAllowedErr)
	}
}

// GetList return user list
// @summary Get list by ID
// @description Gets a list by its ID
// @id get-list
// @produce json
// @Param id path string true "List ID formatted as an UUID string"
// @Success 200 {object} APIResponse
// @Success 400 {object} APIResponse
// @Success 404 {object} APIResponse
// @Success 405 {object} APIResponse
// @Router /api/v1/lists/{id} [get]
// @tags Lists
func (h *APIHandler) GetList(w http.ResponseWriter, r *http.Request) {
	//TODO not implemented yet
	_, err := w.Write([]byte("GetLists not implemented yet"))
	if err != nil {
		h.Log().Error(err)
	}
}

func (h *APIHandler) handleOpenAPIDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, h.apiDoc)
}

// Helpers

func (h *APIHandler) User(r *http.Request) (userID string, err error) {
	// Authentication mechanism not yet established.
	// WIP: A hardcoded value is returned for now.
	uid := "e1263c73-521b-41b5-96e5-58c3f71e65a1"

	ok := uuid.Validate(uid)
	if !ok {
		return "", NoUserErr
	}

	return uid, nil
}

// List returns the list ID from request context.
// Chi router + OpenAPI makes this unnecessary but can be useful when using
// stdlib or a Chi router custom middleware.
func (h *APIHandler) List(r *http.Request) (listID string, err error) {
	val := r.Context().Value(ListCtxKey)
	if val != nil {
		switch v := val.(type) {
		case string:
			return v, nil
		default:
			return listID, InvalidValueTypeErr
		}
	}

	return listID, ListNotFoundErr
}

// closeBody close the body and log errors if happened.
func (h *APIHandler) closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		h.Log().Error(errors.Wrap(err, "failed to close body"))
	}
}

// Handler interface

// Service returns a list svc implementation.
func (h *APIHandler) Service() service.ListService {
	return h.svc
}
