package http

import (
	"fmt"
	"net/http"

	"github.com/vanillazen/stl/backend/internal/sys"
)

type (
	Probe struct {
		sys.Core
	}
)

func NewProbe(name string, opts ...sys.Option) *Probe {
	return &Probe{
		Core: sys.NewCore(name, opts...),
	}
}

func (p Probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "OK")
	if err != nil {
		p.Log().Errorf("%s error %w", p.Name(), err)
	}
}
