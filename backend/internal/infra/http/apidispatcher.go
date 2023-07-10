package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/vanillazen/stl/backend/internal/sys/errors"
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

// ResourceInfo holds information extracted from a URL path. It captures the hierarchical levels and associated IDs
// of a resource based on the segments in the URL path.
//
// For example, given the URL path "api/v1/books/123/chapters/456/paragraphs", the ResourceInfo struct would contain:
// - Levels: ["lists", "tasks"]
// - IDs: ["043a415f-93e3-4f95-97c6-b4d9ca564188"]
// - Error: nil
//
// The Levels field represents the resource levels extracted from the URL path, such as "lists" and "tasks",
// The IDs field stores the associated IDs for each level in the hierarchy, corresponding to the extracted levels.
// The Error field will be nil if the URL path extraction is successful. Otherwise, if an error occurs during the extraction,
// it will contain the corresponding error.
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
			resourceInfo.Error = errors.New("Invalid URL")
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

// Level1 returns the first level of the resource extracted from the URL path.
// If the Levels field has at least one element, Level1 returns that element.
// Otherwise, it returns an empty string.
func (ri ResourceInfo) Level1() (l1 string) {
	if len(ri.Levels) > 0 {
		return ri.Levels[0]
	}

	return l1
}

// IDLevel1 returns the ID associated with the first level of the resource hierarchy.
// If the IDs field has at least one element, IDLevel1 returns that element.
// Otherwise, it returns an empty string.
func (ri ResourceInfo) IDLevel1() (l1id string) {
	if len(ri.IDs) > 0 {
		return ri.IDs[0]
	}

	return l1id
}

// Level2 returns the second level of the resource extracted from the URL path.
// If the Levels field has at least two elements, Level2 returns the second element.
// Otherwise, it returns an empty string.
func (ri ResourceInfo) Level2() (l2 string) {
	if len(ri.Levels) > 1 {
		return ri.Levels[1]
	}

	return l2
}

// IDLevel2 returns the ID associated with the second level of the resource hierarchy.
// If the IDs field has at least two elements, IDLevel2 returns the second element.
// Otherwise, it returns an empty string.
func (ri ResourceInfo) IDLevel2() (l2id string) {
	if len(ri.IDs) > 1 {
		return ri.IDs[1]
	}

	return l2id
}

// Level3 returns the third level of the resource extracted from the URL path.
// If the Levels field has at least three elements, Level3 returns the third element.
// Otherwise, it returns an empty string.
func (ri ResourceInfo) Level3() (l3 string) {
	if len(ri.Levels) > 2 {
		return ri.Levels[2]
	}

	return l3
}

// IDLevel3 returns the ID associated with the third level of the resource hierarchy.
// If the IDs field has at least three elements, IDLevel3 returns the third element.
// Otherwise, it returns an empty string.
func (ri ResourceInfo) IDLevel3() (l3id string) {
	if len(ri.IDs) > 2 {
		return ri.IDs[2]
	}

	return l3id
}
