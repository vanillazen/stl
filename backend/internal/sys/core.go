package sys

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/vanillazen/stl/backend/internal/sys/config"
	"github.com/vanillazen/stl/backend/internal/sys/log"
)

type (
	Core interface {
		Name() string
		Log() log.Logger
		Cfg() *config.Config
		Setup(ctx context.Context) error
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}
)

type (
	SimpleCore struct {
		name     string
		log      log.Logger
		cfg      *config.Config
		didSetup bool
		didStart bool
	}
)

func NewCore(name string, opts ...Option) *SimpleCore {
	name = GenName(name, "worker")

	bw := &SimpleCore{
		name: name,
	}

	for _, opt := range opts {
		opt(bw)
	}

	return bw
}

func (sc *SimpleCore) Name() string {
	return sc.name
}

func (sc *SimpleCore) SetName(name string) {
	sc.name = name
}

func (sc *SimpleCore) Log() log.Logger {
	return sc.log
}

func (sc *SimpleCore) SetLog(log log.Logger) {
	sc.log = log
}

func (sc *SimpleCore) Cfg() *config.Config {
	return sc.cfg
}

func (sc *SimpleCore) SetCfg(cfg *config.Config) {
	sc.cfg = cfg
}

func (sc *SimpleCore) Setup(ctx context.Context) error {
	sc.Log().Infof("%s setup", sc.Name())
	return nil
}

func (sc *SimpleCore) Start(ctx context.Context) error {
	sc.Log().Infof("%s start", sc.Name())
	return nil
}

func (sc *SimpleCore) Stop(ctx context.Context) error {
	sc.Log().Infof("%s stop", sc.Name())
	return nil
}

func GenName(name, defName string) string {
	if strings.Trim(name, " ") == "" {
		return fmt.Sprintf("%s-%s", defName, nameSufix())
	}
	return name
}

func nameSufix() string {
	digest := hash(time.Now().String())
	return digest[len(digest)-8:]
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

type (
	Option func(w *SimpleCore)
)

func WithConfig(cfg *config.Config) Option {
	return func(svc *SimpleCore) {
		svc.SetCfg(cfg)
	}
}

func WithLogger(log log.Logger) Option {
	return func(svc *SimpleCore) {
		svc.SetLog(log)
	}
}
