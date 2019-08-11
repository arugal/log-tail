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
	ServerCnf   = config2.GetDefaultServerConf()
	CommonCnf   = config2.GetDefaultCommonConf()
	CatalogsCnf = config2.GetDefaultCatalogsConf()
)

func Load(cfgFile string) error {
	cnfMap, err := config2.ReaderConfigFromYaml(cfgFile)
	if err != nil {
		return err
	}

	server := cnfMap[Server]
	if server != nil {
		severMap := server.(map[interface{}]interface{})
		ServerCnf, err = config2.UnmarshalServerConfFromYaml(severMap)
		if err != nil {
			return err
		}
	}

	common := cnfMap[Common]
	if common != nil {
		commonMap := common.(map[interface{}]interface{})
		CommonCnf, err = config2.UnmarshalCommonConfFromYaml(commonMap)

		if err != nil {
			return err
		}
	}
	CommonCnf.CfgFile = cfgFile

	catalog := cnfMap[Catalogs]
	if catalog != nil {
		cataLogSlice := catalog.([]interface{})
		CatalogsCnf, err = config2.UnmarshalCatalogConfFromYaml(cataLogSlice)
	}
	return err
}
