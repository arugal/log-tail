package config

import (
	"github.com/vaughan0/go-ini"
	"strings"
)

type CatalogConf struct {
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	IgnoreSuffixs []string `json:"ignore_suffix"`
	ChildFile     []string `json:"child_file"`
}

func (cf *CatalogConf) HasChildFile(childfile string) (ok bool) {
	for _, file := range cf.ChildFile {
		if file == childfile {
			return true
		}
	}
	return false
}

func (cf *CatalogConf) FullFilePath(filePath string) string {
	return cf.Path + "/" + filePath
}

func (cf *CatalogConf) Check() (err error) {
	return nil
}

func LoadAllCatalogFromIni(content string) (catalogs map[string]*CatalogConf, err error) {
	conf, errRet := ini.Load(strings.NewReader(content))
	if errRet != nil {
		err = errRet
		return
	}

	var (
		tmpStr string
		ok     bool
	)
	catalogs = make(map[string]*CatalogConf)
	for name, section := range conf {
		if name == "common" {
			continue
		}

		catalog := CatalogConf{
			Name: name,
		}

		if tmpStr, ok = section["path"]; ok {
			if len(tmpStr) > 1 && strings.HasSuffix(tmpStr, "/") {
				tmpStr = tmpStr[:len(tmpStr)-1]
			}
			catalog.Path = tmpStr
		}

		if tmpStr, ok = section["ignore_suffix"]; ok {
			catalog.IgnoreSuffixs = ParseIgnoreSuffix(tmpStr)
		}
		catalogs[name] = &catalog
	}

	return catalogs, nil
}
