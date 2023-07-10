package http_test

import (
	"testing"

	"github.com/vanillazen/stl/backend/internal/infra/http"
	"github.com/vanillazen/stl/backend/internal/sys/errors"
)

func TestGetResourceInfo(t *testing.T) {
	tests := []struct {
		name           string
		parts          []string
		expectedLevels []string
		expectedIDs    []string
		expectedError  errors.Error
	}{
		{
			name:           "Valid URL parts",
			parts:          []string{"lists", "c5e13593-7903-4f44-9c0b-a6daf28e5763", "tasks", "15da8e3b-ecae-4e63-a721-4851ab0b0b35", "categories", "b12068f8-98eb-46c0-a8b7-62ea3d5e6a99"},
			expectedLevels: []string{"lists", "tasks", "categories"},
			expectedIDs:    []string{"c5e13593-7903-4f44-9c0b-a6daf28e5763", "15da8e3b-ecae-4e63-a721-4851ab0b0b35", "b12068f8-98eb-46c0-a8b7-62ea3d5e6a99"},
			expectedError:  errors.Empty,
		},
		{
			name:           "Invalid URL parts count",
			parts:          []string{"lists"},
			expectedLevels: nil,
			expectedIDs:    nil,
			expectedError:  errors.New("Invalid URL"),
		},
		{
			name:           "Invalid URL ID",
			parts:          []string{"lists", "c5e13593-7903-4f44-9c0b-a6daf28e5763", "tasks", "invalid-id", "categories", "b12068f8-98eb-46c0-a8b7-62ea3d5e6a99"},
			expectedLevels: []string{"lists", "tasks"},
			expectedIDs:    []string{"c5e13593-7903-4f44-9c0b-a6daf28e5763"},
			expectedError:  errors.New("Invalid URL"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := http.GetResourceInfo(test.parts)

			// Verify the levels
			if !equalSlices(result.Levels, test.expectedLevels) {
				t.Errorf("Levels: expected %v, got %v", test.expectedLevels, result.Levels)
			}

			// Verify the IDs
			if !equalSlices(result.IDs, test.expectedIDs) {
				t.Errorf("IDs: expected %v, got %v", test.expectedIDs, result.IDs)
			}

			// Verify the error
			if (result.Error == errors.Empty && test.expectedError != errors.Empty) ||
				(result.Error != errors.Empty && test.expectedError == errors.Empty) ||
				(result.Error != errors.Empty && test.expectedError != errors.Empty &&
					result.Error.Unwrap().Error() != test.expectedError.Unwrap().Error()) {
				t.Errorf("Error: expected '%v', got '%v'", test.expectedError, result.Error)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
