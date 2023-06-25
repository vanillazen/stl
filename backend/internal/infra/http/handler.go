package http

import (
	"io"
	"net/http"

	"github.com/vanillazen/stl/backend/internal/domain/service"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type (
	ListHTTPHandler interface {
		sys.Core
		Service() service.ListService
	}

	ListHandler struct {
		*sys.SimpleCore
		svc service.ListService
	}
)

func NewListHandler(svc service.ListService, opts ...sys.Option) *ListHandler {
	return &ListHandler{
		SimpleCore: sys.NewCore("list-handler", opts...),
		svc:        svc,
	}
}

func (h *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	//TODO not implemented yet
	_, err := w.Write([]byte("GetLists not implemented yet"))
	if err != nil {
		h.Log().Error(err)
	}
}

// Helpers

func (h *ListHandler) User(r *http.Request) (userID string, err error) {
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
func (h *ListHandler) List(r *http.Request) (listID string, err error) {
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
func (h *ListHandler) closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		h.Log().Error(errors.Wrap("failed to close body", err))
	}
}

// Handler interface

// Service returns a list svc implementation.
func (h *ListHandler) Service() service.ListService {
	return h.svc
}
