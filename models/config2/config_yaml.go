package config2

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

type Conf interface {
	Verify() bool
}

type ServerConf struct {
	Host   string      `json:"host"`
	Port   int         `json:"port"`
	Secure *SecureConf `json:"secure"`
}

func (cnf ServerConf) Verify() bool {
	if !cnf.Secure.Verify() {
		return false
	}
	return true
}

type CommonConf struct {
	LastReadOffset        int64         `json:"last_read_offset"`
	ConnMaxTime           int           `json:"conn_max_time"`
	ConnMaxTimeDuration   time.Duration `json:"conn_max_time_duration"`
	HeartInterval         int           `json:"heart_interval"`
	HeartIntervalDuration time.Duration `json:"heart_interval_duration"`
	Log                   *LogConf      `json:"log"`
	Ignore                *IgnoreConf   `json:"ignore"`
	CfgFile               string        `json:"cfg_file"`
}

func (cnf CommonConf) Verify() bool {
	if !cnf.Log.Verify() {
		return false
	}

	if !cnf.Ignore.Verify() {
		return false
	}
	return true
}

func (cnf CommonConf) HeartIntervalFunc() time.Duration {
	return cnf.HeartIntervalDuration * time.Duration(cnf.HeartInterval)
}

func (cnf CommonConf) ConnMaxTimeFunc() time.Duration {
	return cnf.ConnMaxTimeDuration * time.Duration(cnf.ConnMaxTime)
}

type CatalogsConf struct {
	Catalogs []*CatalogConf `json:"catalogs"`
}

func (cnf CatalogsConf) Verify() bool {

	for _, catalog := range cnf.Catalogs {
		if !catalog.Verify() {
			return false
		}
	}

	return true
}

type LogConf struct {
	File    string `json:"file"`
	Way     string `json:"way"`
	Level   string `json:"level"`
	MaxDays int64  `json:"max_days"`
}

func (cnf LogConf) Verify() bool {
	return true
}

type SecureConf struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

func (cnf SecureConf) Verify() bool {
	return true
}

type CatalogConf struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Ignore    *IgnoreConf `json:"ignore"`
	ChildFile []string    `json:"child_file"`
}

func (cnf CatalogConf) Verify() bool {
	if !cnf.Ignore.Verify() {
		return false
	}
	return true
}

func (cnf *CatalogConf) HasChildFile(childfile string) (ok bool) {
	for _, file := range cnf.ChildFile {
		if file == childfile {
			return true
		}
	}
	return false
}

func (cnf *CatalogConf) FullFilePath(filePath string) string {
	return cnf.Path + "/" + filePath
}

type IgnoreConf struct {
	Suffix []string `json:"suffix"`
	Regexp []string `json:"regexp"`
}

func (cnf IgnoreConf) Verify() bool {
	return true
}

func GetDefaultServerConf() *ServerConf {
	conf := &ServerConf{
		Host: "0.0.0.0",
		Port: 3000,
		Secure: &SecureConf{
			User: "admin",
			Pwd:  "admin",
		},
	}
	return conf
}

func GetDefaultCommonConf() *CommonConf {
	conf := &CommonConf{
		LastReadOffset:        int64(1000),
		ConnMaxTime:           10,
		ConnMaxTimeDuration:   time.Minute,
		HeartInterval:         10,
		HeartIntervalDuration: time.Second,
		Log: &LogConf{
			File:    "console",
			Way:     "console",
			Level:   "info",
			MaxDays: int64(7),
		},
		Ignore: &IgnoreConf{
			Suffix: []string{".jar", "war", ".html", ".js", ".css", ".java", ".class", ".gz", ".tar", ".zip", ".rar", ".jgp", ".png", ".xls", ".xlxs", ".pdf"},
			Regexp: []string{},
		},
	}
	return conf
}

func GetDefaultCatalogsConf() *CatalogsConf {
	conf := &CatalogsConf{
		Catalogs: []*CatalogConf{},
	}
	return conf
}

func GetDefaultCatalogConf() *CatalogConf {
	conf := &CatalogConf{
		Ignore: &IgnoreConf{},
	}
	return conf
}

func UnmarshalServerConfFromYaml(confMap map[interface{}]interface{}) (*ServerConf, error) {
	cfg := GetDefaultServerConf()
	err := fillConfFromMap(confMap, cfg)

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func UnmarshalCommonConfFromYaml(confMap map[interface{}]interface{}) (*CommonConf, error) {
	cfg := GetDefaultCommonConf()
	err := fillConfFromMap(confMap, cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func UnmarshalCatalogConfFromYaml(confSlice []interface{}) (*CatalogsConf, error) {
	cfg := GetDefaultCatalogsConf()

	err := fillConfFromSlice(confSlice, cfg)

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func fillConfFromSlice(confSlice []interface{}, conf *CatalogsConf) error {
	for _, confMap := range confSlice {
		childCnf := GetDefaultCatalogConf()
		err := fillConfFromMap(confMap.(map[interface{}]interface{}), childCnf)
		if err != nil {
			return err
		}
		conf.Catalogs = append(conf.Catalogs, childCnf)
	}
	return nil
}

func fillConfFromMap(confMap map[interface{}]interface{}, conf Conf) error {
	for k, v := range confMap {
		if v == nil {
			continue
		}
		err := setField(conf, k.(string), v)
		if err != nil {
			return err
		}
	}

	if !conf.Verify() {
		return fmt.Errorf("verify failed %v", conf)
	}

	return nil
}

func setField(conf Conf, name string, value interface{}) error {
	structValue := reflect.ValueOf(conf).Elem()
	name = strings.ReplaceAll(name, "_", "")
	structFieldValue := structValue.FieldByNameFunc(func(field string) bool {
		if strings.ToLower(field) == name {
			return true
		}
		return false
	})

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such filed: %s in conf", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	if structFieldType != reflect.TypeOf(value) {
		var err error
		val, err = typeConversion(value, structFieldValue.Interface())
		if err != nil {
			return err
		}
	}
	structFieldValue.Set(val)
	return nil
}

func typeConversion(value interface{}, rowValue interface{}) (reflect.Value, error) {
	switch rowValue.(type) {
	case int8:
		return reflect.ValueOf(int8(value.(int))), nil
	case int32:
		return reflect.ValueOf(int32(value.(int))), nil
	case int64:
		return reflect.ValueOf(int64(value.(int))), nil
	case []string:
		value := value.([]interface{})
		var strValue []string
		for _, v := range value {
			strValue = append(strValue, v.(string))
		}
		return reflect.ValueOf(strValue), nil
	default:
		err := fillConfFromMap(value.(map[interface{}]interface{}), rowValue.(Conf))
		return reflect.ValueOf(rowValue), err
	}
}

func ReaderConfigFromYaml(path string) (configMap map[string]interface{}, err error) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(buffer, &configMap)
	return
}
