package ci

import (
	"fmt"
	"log-tail/models/config"
	"testing"
)

func TestParseServerCfg(t *testing.T) {
	cfg, err := config.UnmarshalServerConfFromIni(nil, Content)
	if err != nil {
		t.Error(err)
		t.Fail()
	} else {
		t.Log(cfg)
	}
}

func TestParseCatalogCfg(t *testing.T) {
	catalogs, err := config.LoadAllCatalogFromIni(Content)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(catalogs)
}
