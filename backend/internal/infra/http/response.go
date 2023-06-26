package http

import (
	"encoding/json"
	"net/http"

	"github.com/vanillazen/stl/backend/internal/sys/errors"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Count   int         `json:"count,omitempty"`
	Pages   int         `json:"pages,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *ListHandler) handleSuccess(w http.ResponseWriter, payload interface{}, count, pages int, msg ...string) {
	var m string
	if len(msg) > 0 {
		m = msg[0]
	}

	response := APIResponse{
		Success: true,
		Message: m,
		Count:   count,
		Pages:   pages,
		Data:    payload,
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Log().Error(errors.Wrap(err, "error encoding handler success"))
	}

	return
}

func (h *ListHandler) handleError(w http.ResponseWriter, handlerError error) {
	response := APIResponse{
		Success: false,
		Message: handlerError.Error(),
	}

	h.Log().Error("handler error:", handlerError)

	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Log().Error(errors.Wrap(err, "error encoding handler error"))
	}

	return
}
