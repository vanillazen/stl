package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vanillazen/stl/backend/internal/domain/service"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/transport"
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
		h.handleError(w, 0, MethodNotAllowedErr)
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
	ctx := r.Context()

	userID, err := h.User(r)
	if err != nil {
		h.handleError(w, http.StatusNotFound, errors.Wrap(err))
		return
	}

	resource, ok := h.resourceID(r)
	if !ok {
		h.handleError(w, 0, errors.Wrap(NoResourceErr, "get list error"))
		return
	}

	req := transport.GetListReq{
		UserID: userID,
		ListID: resource.L1ID(),
	}

	res := h.Service().GetList(ctx, req)
	if err = res.Err(); err != nil {
		err = errors.Wrap(err, "get list error")
		h.handleError(w, 0, err)
		return
	}

	h.handleSuccess(w, res, 1, 1)
}

func (h *APIHandler) handleTask(w http.ResponseWriter, r *http.Request) {
	h.Log().Error("not implemented yet")
}

func (h *APIHandler) handleOpenAPIDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, h.apiDoc)
}

// Helpers

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
