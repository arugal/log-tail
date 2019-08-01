package g

import (
	"github.com/arugal/log-tail/models/config2"
)

const (
	Server   = "server"
	Common   = "common"
	Catalogs = "catalogs"
)

var (
	ServerCnf   *config2.ServerConf
	CommonCnf   *config2.CommonConf
	CatalogsCnf *config2.CatalogsConf
)

func Load(cfgFile string) error {
	cnfMap, err := config2.ReaderConfigFromYaml(cfgFile)
	if err != nil {
		return err
	}

	server := cnfMap[Server]
	severMap := server.(map[interface{}]interface{})
	ServerCnf, err = config2.UnmarshalServerConfFromYaml(severMap)
	if err != nil {
		return err
	}

	common := cnfMap[Common]
	commonMap := common.(map[interface{}]interface{})
	CommonCnf, err = config2.UnmarshalCommonConfFromYaml(commonMap)

	if err != nil {
		return err
	}
	CommonCnf.CfgFile = cfgFile

	catalog := cnfMap[Catalogs]
	cataLogSlice := catalog.([]interface{})
	CatalogsCnf, err = config2.UnmarshalCatalogConfFromYaml(cataLogSlice)
	return err
}
