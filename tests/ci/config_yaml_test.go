package ci

import (
	"github.com/arugal/log-tail/g"
	"github.com/arugal/log-tail/models/config2"
	"testing"
)

var confMap map[string]interface{}

func init() {
	var err error
	confMap, err = config2.ReaderConfigFromYaml("./config.yaml")
	if err != nil {
		panic(err)
	}
}

func assert(conf interface{}, err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}

func TestServerYamlParse(t *testing.T) {
	server := confMap[g.Server]
	severMap := server.(map[interface{}]interface{})
	serverCfg, err := config2.UnmarshalServerConfFromYaml(severMap)
	assert(serverCfg, err, t)
}

func TestCommonYamlParse(t *testing.T) {
	common := confMap[g.Common]
	commonMap := common.(map[interface{}]interface{})
	commonCfg, err := config2.UnmarshalCommonConfFromYaml(commonMap)
	assert(commonCfg, err, t)
}

func TestCatalogYamlParse(t *testing.T) {
	catalog := confMap[g.Catalogs]
	cataLogSlice := catalog.([]interface{})
	catalogConfs, err := config2.UnmarshalCatalogConfFromYaml(cataLogSlice)
	assert(catalogConfs, err, t)
}
