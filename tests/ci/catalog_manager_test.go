package ci

import (
	"fmt"
	"log-tail/server/catalog"
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

func TestGo(t *testing.T) {

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		fmt.Println("1:" + <-c1)
	}()

	go func() {
		fmt.Println("2:" + <-c2)
	}()

	c1 <- "abc"
	c2 <- "dfg"

	time.Sleep(time.Second)

}
