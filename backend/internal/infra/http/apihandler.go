package http

import (
	"encoding/json"
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
	res, ok := h.resource(r)
	if !ok {
		h.handleError(w, http.StatusBadRequest, NoResourceErr)
	}

	switch r.Method {
	case http.MethodGet:
		if res.IDLevel1() != "" {
			h.GetList(w, r)
			return
		} else {
			h.GetLists(w, r)
			return
		}

	case http.MethodPost:
		//h.CreateListList(w, r)
		h.Log().Error("not implemented yet")

	case http.MethodPut:
		//h.UpdateList(w, r)
		h.Log().Error("not implemented yet")

	default:
		h.handleError(w, http.StatusMethodNotAllowed, MethodNotAllowedErr)
	}
}

func (h *APIHandler) GetLists(w http.ResponseWriter, r *http.Request) {
	h.Log().Error("not implemented yet")
}

// CreateList creates a new list
// @summary Create a new list
// @description Creates a new list with the provided details
// @id create-list
// @accept json
// @produce json
// @Param list body transport.CreateListReq true "List name and description"
// @Success 201 {object} APIResponse
// @Success 400 {object} APIResponse
// @Router /api/v1/lists [post]
// @tags Lists
func (h *APIHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.User(r)
	if err != nil {
		h.handleError(w, http.StatusNotFound, errors.Wrap(err))
		return
	}

	var req transport.CreateListReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(w, http.StatusBadRequest, errors.Wrap(err, "invalid request payload"))
		return
	}

	req.UserID = userID

	res := h.Service().CreateList(ctx, req)
	if err = res.Err(); err != nil {
		err = errors.Wrap(err, "get list error")
		h.handleError(w, http.StatusNotFound, err)
		return
	}

	h.handleSuccess(w, res, 1, 1)

	//h.Service().CreateList(ctx, list)
	//
	//h.handleSuccess(w, list, 1, 1)
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

	resource, ok := h.resource(r)
	if !ok {
		h.handleError(w, http.StatusBadRequest, errors.Wrap(NoResourceErr, "get list error"))
		return
	}

	req := transport.GetListReq{
		UserID: userID,
		ListID: resource.IDLevel1(),
	}

	res := h.Service().GetList(ctx, req)
	if err = res.Err(); err != nil {
		err = errors.Wrap(err, "get list error")
		h.handleError(w, http.StatusNotFound, err)
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
