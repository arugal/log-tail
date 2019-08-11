package server

import (
	"fmt"
	"github.com/arugal/log-tail/assets"
	"github.com/arugal/log-tail/g"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"

	tailNet "github.com/arugal/log-tail/util/net"
)

var (
	httpServerReadTimeout  = 10 * time.Second
	httpServerWriteTimeout = 10 * time.Second
)

func (srv *Service) RunDashboardServer() (err error) {
	serverCnf := g.ServerCnf

	// load static resource
	err = assets.Load("")
	if err != nil {
		return err
	}

	// url router
	router := mux.NewRouter()
	router.Use(tailNet.NewHttpAuthMiddleware(serverCnf.Secure.User, serverCnf.Secure.Pwd).Middleware)
	router.Use(tailNet.NewCrossDomainMiddleware().Middleware)

	// api
	router.HandleFunc("/api/catalog", srv.GetCataLogInfo).Methods("GET")
	router.HandleFunc("/api/tail/{catalog}/{file}", srv.GetLogTail)

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

	address := fmt.Sprintf("%s:%d", serverCnf.Host, serverCnf.Port)
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

	srv.log.Info("Start Dashboard listen %s", ln.Addr())
	server.Serve(ln)
	return nil
}
