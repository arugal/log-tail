package ci

import (
	"github.com/Arugal/log-tail/g"
	"github.com/Arugal/log-tail/models/config"
)

var (
	Content string
)

func init() {
	g.GlbServerCfg.CfgFile = "parse_test.ini"
	content, err := config.GetRenderedConfFromFile(g.GlbServerCfg.CfgFile)
	if err != nil {
		panic(err)
	}
	cfg, err := config.UnmarshalServerConfFromIni(nil, content)
	if err != nil {
		panic(err)
	}
	g.GlbServerCfg.ServerCommonConf = *cfg
	Content = content
}
