package ci

import (
	"fmt"
	"github.com/Arugal/log-tail/models/config"
	"regexp"
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

func TestRegexp(t *testing.T) {
	reg := ".log."
	match, _ := regexp.MatchString(reg, "application.log.2019-05-16")
	fmt.Println(match)
}
