package server

import (
	"fmt"
	"github.com/Arugal/log-tail/assets"
	"github.com/Arugal/log-tail/g"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"

	tailNet "github.com/Arugal/log-tail/util/net"
)

var (
	httpServerReadTimeout  = 10 * time.Second
	httpServerWriteTimeout = 10 * time.Second
)

func (svr *Service) RunDashboardServer(addr string, port int) (err error) {
	// url router
	router := mux.NewRouter()
	router.Use(tailNet.NewHttpAuthMiddleware(g.GlbServerCfg.User, g.GlbServerCfg.Pwd).Middleware)
	router.Use(tailNet.NewCrossDomainMiddleware().Middleware)

	// api
	router.HandleFunc("/api/catalog", svr.GetCataLogInfo).Methods("GET")
	router.HandleFunc("/api/tail/{catalog}/{file}", svr.GetLogTail)

	// view
	router.Handle("/favicon.ico", http.FileServer(assets.FileSystem)).Methods("GET")
	router.PathPrefix("/static/").Handler(tailNet.MakeHttpGzipHandler(http.StripPrefix("/static/",
		http.FileServer(assets.FileSystem)))).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static", http.StatusMovedPermanently)
	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
	})

	address := fmt.Sprintf("%s:%d", addr, port)
	if address == "" {
		address = ":3000"
	}

	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  httpServerReadTimeout,
		WriteTimeout: httpServerWriteTimeout,
	}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	svr.log.Info("Start Dashboard listen %s", ln.Addr())
	server.Serve(ln)
	return nil
}
