package server

import (
	"github.com/arugal/log-tail/server/catalog"
	"github.com/arugal/log-tail/server/control"
	"github.com/arugal/log-tail/util/log"
)

type Service struct {
	cm  *catalog.CatalogManger
	cm2 *control.ConnManager
	log log.Logger
}

func NewService() (srv *Service, err error) {
	cm, err := catalog.NewCataLogManager()
	if err != nil {
		return nil, err
	}
	srv = &Service{
		cm:  cm,
		cm2: control.NewConnManager(),
		log: log.NewPrefixLogger("service"),
	}

	return srv, nil
}

func (srv *Service) Start() {
	go srv.cm.Run()
	go srv.cm2.Run()
	srv.RunDashboardServer()
}
