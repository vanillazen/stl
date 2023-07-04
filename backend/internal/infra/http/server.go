package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vanillazen/stl/backend/internal/domain/service"
	"github.com/vanillazen/stl/backend/internal/sys"
	"github.com/vanillazen/stl/backend/internal/sys/config"
)

type (
	Server struct {
		sys.Core
		opts []sys.Option
		http.Server
		*ServeMux
		apiV1 *APIHandler
		svc   service.ListService
	}
)

const (
	apiV1          = "/api/v1/"
	apiV1Docs      = apiV1 + "docs/"
	httpServerName = "http-server"
)

var (
	cfgKey = config.Key
)

func NewServer(svc service.ListService, apiDoc string, opts ...sys.Option) (server *Server) {
	apiHandler := NewAPIHandler(svc, apiDoc, opts...)

	return &Server{
		Core:     sys.NewCore("api-server", opts...),
		opts:     opts,
		ServeMux: NewServeMux("api-router", opts...),
		apiV1:    apiHandler,
		svc:      svc,
	}
}

func (srv *Server) Setup(ctx context.Context) error {
	//reqLog := NewReqLoggerMiddleware(srv.Log())

	// TODO: Add middlewares for srv.router:
	// RequestID, RealIP, Logging and Recover

	// TODO: Setup Mux routes & handlers
	srv.Mux().HandleFunc(apiV1Docs, srv.apiV1.handleOpenAPIDocs)
	srv.Mux().HandleFunc(apiV1, srv.apiV1.handleV1)

	return nil
}

func (srv *Server) Start(ctx context.Context) error {
	srv.Server = http.Server{
		Addr:    srv.Address(),
		Handler: srv.Mux(),
	}

	var group, errGrpCtx = errgroup.WithContext(ctx)
	group.Go(func() error {
		srv.Log().Infof("%s listening at %s", srv.Name(), srv.Address())
		defer srv.Log().Errorf("%s shutdown", srv.Name())

		err := srv.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-errGrpCtx.Done()
		srv.Log().Errorf("%s shutdown", srv.Name())

		ctx, cancel := context.WithTimeout(context.Background(), srv.ShutdownTimeout())
		defer cancel()

		if err := srv.Server.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}

func (srv *Server) SetMux(sm *ServeMux) {
	srv.ServeMux = sm
}

func (srv *Server) Mux() (m *ServeMux) {
	return srv.ServeMux
}

func (srv *Server) Mount(pattern string, handler http.Handler) {
	srv.ServeMux.Mount(pattern, handler)
}

func (srv *Server) Address() string {
	host := srv.Cfg().GetString(cfgKey.APIServerHost)
	port := srv.Cfg().GetInt(cfgKey.APIServerPort)
	return fmt.Sprintf("%s:%d", host, port)
}

func (srv *Server) ShutdownTimeout() time.Duration {
	secs := time.Duration(srv.Cfg().GetInt(cfgKey.APIServerTimeout))
	return secs * time.Second
}
