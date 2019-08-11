package catalog

import (
	"fmt"
	"github.com/arugal/log-tail/g"
	"github.com/arugal/log-tail/models/config2"
	"github.com/arugal/log-tail/util/log"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type CatalogManger struct {
	catalogs map[string]*config2.CatalogConf
	log      log.Logger
}

func NewCataLogManager() (cm *CatalogManger, err error) {
	cm = &CatalogManger{
		catalogs: map[string]*config2.CatalogConf{},
		log:      log.NewPrefixLogger("catalog-manager"),
	}

	ok := cm.AddCatalogs(g.CatalogsCnf.Catalogs)
	if !ok {
		return nil, fmt.Errorf("init CatalogMager failure, please invalid cfgFile %s and see the log %s",
			g.CommonCnf.CfgFile, g.CommonCnf.Log.File)
	}
	return cm, nil
}

func (m *CatalogManger) Run() {

	go func(m *CatalogManger) {
		m.log.Info("Start CatalogManger  scanInterval %s", time.Duration(time.Minute*5).String())
		timer := time.NewTicker(time.Minute * 5)
		for {
			select {
			case <-timer.C:
				m.refreshAllCatalog()
			}
		}
	}(m)
}

func (m *CatalogManger) GetCatalogInfo(catalog string) (conf config2.CatalogConf, ok bool) {
	realCf, ok := m.catalogs[catalog]
	conf = *realCf
	return
}

func (m *CatalogManger) GetAllCatalogInfo() map[string]*config2.CatalogConf {
	return m.catalogs
}

func (m *CatalogManger) refreshAllCatalog() {
	for _, cnf := range m.catalogs {
		childFile, err := m.refreshCatalog(cnf)
		if err != nil {
			log.Warn("interval refresh catalog failure [%s] : %s case:%v", cnf.Name, cnf.Path, err)
		} else {
			cnf.ChildFile = childFile
		}
	}
}

func (m *CatalogManger) AddCatalogs(catalogs []*config2.CatalogConf) (ok bool) {
	for _, cnf := range catalogs {
		childFile, err := m.refreshCatalog(cnf)
		if err != nil {
			log.Warn("refresh catalog failure [%s] : %s case:%v", cnf.Name, cnf.Path, err)
			return false
		}
		cnf.ChildFile = childFile
		m.catalogs[cnf.Name] = cnf
	}
	return true
}

func (m *CatalogManger) AddCataLog(cnf *config2.CatalogConf) (ok bool) {
	childFile, err := m.refreshCatalog(cnf)
	if err != nil {
		log.Warn("refresh catalog failure [%s] : %s case:%v", cnf.Name, cnf.Path, err)
		return false
	}
	cnf.ChildFile = childFile
	m.catalogs[cnf.Name] = cnf
	return true
}

func (m *CatalogManger) refreshCatalog(cnf *config2.CatalogConf) (childFile []string, err error) {
	log.Debug("refresh catalog [%s] : %v", cnf.Name, cnf)
	pathInfo, err := os.Stat(cnf.Path)
	if err != nil {
		return nil, err
	}

	if !pathInfo.IsDir() {
		return nil, &UnDirectoryError{Path: cnf.Path}
	}

	childFile = []string{}
	fileInfos, _ := ioutil.ReadDir(cnf.Path)
	for _, fileInfo := range fileInfos {
		var ignore = false
		for _, suffix := range cnf.Ignore.Suffix {
			if strings.HasSuffix(fileInfo.Name(), suffix) {
				ignore = true
				break
			}
		}

		if !ignore {
			for _, reg := range cnf.Ignore.Regexp {
				match, _ := regexp.MatchString(reg, fileInfo.Name())
				if match {
					ignore = true
					break
				}
			}
		}

		if !ignore {
			for _, suffix := range g.CommonCnf.Ignore.Suffix {
				if strings.HasSuffix(fileInfo.Name(), suffix) {
					ignore = true
					break
				}
			}
		}

		if !ignore {
			for _, reg := range g.CommonCnf.Ignore.Regexp {
				match, _ := regexp.MatchString(reg, fileInfo.Name())
				if match {
					ignore = true
					break
				}
			}
		}

		if ignore {
			log.Debug("ignore suffix [%s] : %s", cnf.Name, fileInfo.Name())
			continue
		}

		fullPath := cnf.Path + "/" + fileInfo.Name()
		// file
		if fileInfo.IsDir() {
			log.Debug("ignore is folder [%s] : %s", cnf.Name, fileInfo.Name())
			continue
		}

		// read
		file, err := os.OpenFile(fullPath, os.O_RDONLY, 0666)
		if err != nil {
			log.Debug("ignore dot not read [%s] : %s case:%v", cnf.Name, fileInfo.Name(), err)
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
