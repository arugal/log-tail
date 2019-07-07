package ci

import (
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
	_, err := config.LoadAllCatalogFromIni(Content)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestRegexp(t *testing.T) {
	reg := ".log."
	match, _ := regexp.MatchString(reg, "application.log.2019-05-16")
	if !match {
		t.Fail()
	}
}
