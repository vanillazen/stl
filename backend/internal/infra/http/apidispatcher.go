package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type ResourceInfo struct {
	Levels []string
	IDs    []string
	Error  error
}

type HandlerFunc func(*APIHandler, http.ResponseWriter, *http.Request)

var handlers = map[string]HandlerFunc{
	"lists": (*APIHandler).handleList,
	"task":  (*APIHandler).handleTask,
}

func (h *APIHandler) handleV1(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 || parts[1] != "api" || parts[2] != "v1" {
		msg := "invalid URL"
		h.handleError(w, http.StatusBadRequest, errors.New(msg))
		return
	}

	resourceInfo := getResourceInfo(parts[3:])

	if resourceInfo.Error != nil {
		http.Error(w, resourceInfo.Error.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), ResourceCtxKey, resourceInfo)
	r = r.WithContext(ctx)

	for i := range resourceInfo.Levels {
		handler, ok := handlers[resourceInfo.Levels[i]]
		if !ok {
			http.Error(w, "Invalid resource", http.StatusNotFound)
			return
		}
		handler(h, w, r)
	}
}

func getResourceInfo(parts []string) ResourceInfo {
	resourceInfo := ResourceInfo{}
	levelsCount := len(parts)

	if levelsCount < 2 || levelsCount%2 != 0 {
		resourceInfo.Error = errors.New("Invalid URL")
		return resourceInfo
	}

	for i := 0; i < levelsCount; i += 2 {
		resourceInfo.Levels = append(resourceInfo.Levels, parts[i])
		if i > 0 && !isValidID(parts[i+1]) {
			resourceInfo.Error = fmt.Errorf("Invalid URL")
			return resourceInfo
		}
		resourceInfo.IDs = append(resourceInfo.IDs, parts[i+1])
	}

	return resourceInfo
}

func isValidID(id string) bool {
	return uuid.Validate(id)
}

// Helpers

func (h *APIHandler) User(r *http.Request) (userID string, err error) {
	// Authentication mechanism not yet established.
	// WIP: A hardcoded value is returned for now.
	uid := "0792b97b-4f88-42a8-a035-1d0aad0ae7f8"

	ok := uuid.Validate(uid)
	if !ok {
		return "", NoUserErr
	}

	return uid, nil
}

func (ri ResourceInfo) L1() (l1 string) {
	if len(ri.Levels) > 0 {
		return ri.Levels[0]
	}

	return l1
}

func (ri ResourceInfo) L1ID() (l1id string) {
	if len(ri.IDs) > 0 {
		return ri.IDs[0]
	}

	return l1id
}

func (ri ResourceInfo) L2() (l2 string) {
	if len(ri.Levels) > 1 {
		return ri.Levels[1]
	}

	return l2
}

func (ri ResourceInfo) L2ID() (l2id string) {
	if len(ri.IDs) > 1 {
		return ri.IDs[1]
	}

	return l2id
}

func (ri ResourceInfo) L3() (l3 string) {
	if len(ri.Levels) > 2 {
		return ri.Levels[2]
	}

	return l3
}

func (ri ResourceInfo) L3ID() (l2id string) {
	if len(ri.IDs) > 2 {
		return ri.IDs[2]
	}

	return l2id
}
