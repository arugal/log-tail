package g

import (
	"github.com/Arugal/log-tail/models/config"
)

var (
	GlbServerCfg *ServerCfg
)

func init() {
	GlbServerCfg = &ServerCfg{
		ServerCommonConf: *config.GetDefaultServerConf(),
	}
}

type ServerCfg struct {
	config.ServerCommonConf
	CfgFile string
}
