package http

import (
	"net/http"

	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	ServeMux struct {
		sys.Core
		*http.ServeMux
	}
)

func NewServeMux(name string, opts ...sys.Option) *ServeMux {
	return &ServeMux{
		Core:     sys.NewCore(name, opts...),
		ServeMux: http.NewServeMux(),
	}
}

func (sm *ServeMux) Mount(pattern string, handler http.Handler) {
	// TODO: Implement route-handler mount mechanism
}
