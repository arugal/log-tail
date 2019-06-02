package catalog

import (
	"fmt"
	"io/ioutil"
	"log-tail/g"
	"log-tail/models/config"
	"log-tail/util/log"
	"os"
	"strings"
	"time"
)

type CatalogManger struct {
	catalogs map[string]*config.CatalogConf
	log      log.Logger
}

func NewCataLogManager() (cm *CatalogManger, err error) {
	cm = &CatalogManger{
		catalogs: map[string]*config.CatalogConf{},
		log:      log.NewPrefixLogger("catalog-manager"),
	}

	content, err := config.GetRenderedConfFromFile(g.GlbServerCfg.CfgFile)
	if err != nil {
		return nil, err
	}

	catalogs, err := config.LoadAllCatalogFromIni(content)
	if err != nil {
		return nil, err
	}

	ok := cm.AddCatalogs(catalogs)
	if !ok {
		return nil, fmt.Errorf("init CatalogMager failure, please invalid cfgFile %s and see the log %s",
			g.GlbServerCfg.CfgFile, g.GlbServerCfg.LogFile)
	}
	return cm, nil
}

func (m *CatalogManger) Run() {

	go func(m *CatalogManger) {
		timer := time.NewTimer(time.Minute * 5)
		for {
			select {
			case <-timer.C:
				m.refreshAllCatalog()
			}
		}
	}(m)
}

func (m *CatalogManger) GetCatalogInfo(catalog string) (conf config.CatalogConf, ok bool) {
	realCf, ok := m.catalogs[catalog]
	conf = *realCf
	return
}

func (m *CatalogManger) GetAllCatalogInfo() map[string]*config.CatalogConf {
	return m.catalogs
}

func (m *CatalogManger) refreshAllCatalog() {
	for name, conf := range m.catalogs {
		childFile, err := m.refreshCatalog(name, conf)
		if err != nil {
			log.Warn("interval refresh catalog failure [%s] : %s case:%v", name, conf.Path, err)
		} else {
			conf.ChildFile = childFile
		}
	}
}

func (m *CatalogManger) AddCatalogs(catalogs map[string]*config.CatalogConf) (ok bool) {
	for name, conf := range catalogs {
		childFile, err := m.refreshCatalog(name, conf)
		if err != nil {
			log.Warn("refresh catalog failure [%s] : %s case:%v", conf.Name, conf.Path, err)
			return false
		}
		conf.ChildFile = childFile
		m.catalogs[name] = conf
	}
	return true
}

func (m *CatalogManger) AddCataLog(name string, conf *config.CatalogConf) (ok bool) {
	childFile, err := m.refreshCatalog(name, conf)
	if err != nil {
		log.Warn("refresh catalog failure [%s] : %s case:%v", conf.Name, conf.Path, err)
		return false
	}
	conf.ChildFile = childFile
	m.catalogs[name] = conf
	return true
}

func (m *CatalogManger) refreshCatalog(name string, conf *config.CatalogConf) (childFile []string, err error) {
	log.Debug("refresh catalog [%s] : %v", name, conf)
	pathInfo, err := os.Stat(conf.Path)
	if err != nil {
		return nil, err
	}

	if !pathInfo.IsDir() {
		return nil, &UnDirectoryError{Path: conf.Path}
	}

	childFile = []string{}
	fileInfos, _ := ioutil.ReadDir(conf.Path)
	for _, fileInfo := range fileInfos {
		var ignore = false
		for _, suffix := range conf.IgnoreSuffixs {
			if strings.HasSuffix(fileInfo.Name(), suffix) {
				ignore = true
				break
			}
		}

		if !ignore {
			for _, suffix := range g.GlbServerCfg.IgnoreSuffix {
				if strings.HasSuffix(fileInfo.Name(), suffix) {
					ignore = true
					break
				}
			}
		}

		if ignore {
			log.Debug("ignore suffix [%s] : %s", conf.Name, fileInfo.Name())
			continue
		}

		fullPath := conf.Path + "/" + fileInfo.Name()
		// file
		if fileInfo.IsDir() {
			log.Debug("ignore is folder [%s] : %s", conf.Name, fileInfo.Name())
			continue
		}

		// read
		file, err := os.OpenFile(fullPath, os.O_RDONLY, 0666)
		if err != nil {
			log.Debug("ignore dot not read [%s] : %s case:%v", conf.Name, fileInfo.Name(), err)
			continue
		} else {
			_ = file.Close()
		}
		childFile = append(childFile, fileInfo.Name())
	}
	return childFile, nil
}

type UnDirectoryError struct {
	Path string `json:"path"`
}

func IsUnDirectoryError(err error) bool {
	switch err.(type) {
	case *UnDirectoryError:
		return true
	default:
		return false
	}
}

func (e *UnDirectoryError) Error() string {
	return fmt.Sprintf("path=%s", e.Path)
}
