package server

import (
	"fmt"
	"github.com/Arugal/log-tail/assets"
	"github.com/Arugal/log-tail/g"
	"github.com/Arugal/log-tail/server/catalog"
	"github.com/Arugal/log-tail/server/control"
	"github.com/Arugal/log-tail/util/log"
)

var ServerService *Service

type Service struct {
	cm  *catalog.CatalogManger
	cm2 *control.ConnManager
	log log.Logger
}

func NewService() (srv *Service, err error) {
	cfg := g.GlbServerCfg.ServerCommonConf
	cm, err := catalog.NewCataLogManager()
	if err != nil {
		return nil, err
	}
	srv = &Service{
		cm:  cm,
		cm2: control.NewConnManager(),
		log: log.NewPrefixLogger("service"),
	}

	err = assets.Load(cfg.AssetsDir)
	if err != nil {
		err = fmt.Errorf("Load assets error: %v", err)
		return
	}

	return srv, nil
}

func (srv *Service) Start() {
	go srv.cm.Run()
	go srv.cm2.Run()
	go srv.RunDashboardServer(g.GlbServerCfg.BindAddr, g.GlbServerCfg.BindPort)
}
