package http

import (
	"encoding/json"
	"net/http"

	"github.com/vanillazen/stl/backend/internal/sys/config"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
)

type (
	APIResponse struct {
		Count  int         `json:"count,omitempty"`
		Pages  int         `json:"pages,omitempty"`
		Data   interface{} `json:"data,omitempty"`
		Status `json:"error,omitempty"`
	}

	Status struct {
		OK          bool
		Message     string `json:"message,omitempty"`
		InternalErr string `json:"internalError,omitempty"`
	}
)

func (h *APIHandler) handleSuccess(w http.ResponseWriter, payload interface{}, count, pages int, msg ...string) {
	var m string
	if len(msg) > 0 {
		m = msg[0]
	}

	response := APIResponse{
		Count: count,
		Pages: pages,
		Data:  payload,
		Status: Status{
			OK:      true,
			Message: m,
		},
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Log().Error(errors.Wrap(err, "error encoding handler success"))
	}

	return
}

func (h *APIHandler) handleError(w http.ResponseWriter, httpStatus int, handlerError error, message ...string) {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	}

	var intErr string
	if h.Cfg().GetBool(config.Key.APIErrorExposeInt) {
		intErr = handlerError.Error()
	}

	response := APIResponse{
		Status: Status{
			OK:          false,
			Message:     msg,
			InternalErr: intErr,
		},
	}

	h.Log().Errorf("handler error:\n%s", errors.Stacktrace(handlerError))

	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Log().Error(errors.Wrap(err, "error encoding handler error"))
	}

	return
}
