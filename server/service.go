package server

import (
	"fmt"
	"log-tail/assets"
	"log-tail/g"
	"log-tail/server/catalog"
	"log-tail/server/control"
	"log-tail/util/log"
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
