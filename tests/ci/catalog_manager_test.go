package ci

import (
	"fmt"
	"github.com/arugal/log-tail/server/catalog"
	"testing"
	"time"
)

func TestCatalogManager(t *testing.T) {
	cm, err := catalog.NewCataLogManager()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	cm.Run()

	timer := time.NewTimer(time.Second * 20)
	<-timer.C

	for _, conf := range cm.GetAllCatalogInfo() {
		fmt.Println(conf)
	}
}
