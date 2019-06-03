package config

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vaughan0/go-ini"
)

var (
	glbEnvs map[string]string
)

type ServerCommonConf struct {
	BindAddr string `json:"bind_addr"`
	BindPort int    `json:"bind_port"`
	// time unit minute
	ConnMaxTime time.Duration `json:"con_max_time"`
	// time unit second
	HeartInterval  time.Duration `json:"heart_interval"`
	LogFile        string        `json:"log_file"`
	LogLevel       string        `json:"log_level"`
	LogWay         string        `json:"log_way"`
	LogMaxDays     int64         `json:"log_max_days"`
	User           string        `json:"user"`
	Pwd            string        `json:"pwd"`
	IgnoreSuffix   []string      `json:"ignore_suffix"`
	LastReadOffset int64         `json:"last_read_offset"`
	AssetsDir      string        `json:"assets_dir"`
}

func (cfg *ServerCommonConf) Check() (err error) {
	return nil
}

func GetDefaultServerConf() *ServerCommonConf {
	cfg := &ServerCommonConf{
		BindAddr:       "0.0.0.0",
		BindPort:       3000,
		ConnMaxTime:    10 * time.Minute,
		HeartInterval:  10 * time.Second,
		LogFile:        "console",
		LogWay:         "console",
		LogLevel:       "info",
		LogMaxDays:     int64(64),
		User:           "admin",
		Pwd:            "admin",
		IgnoreSuffix:   []string{},
		LastReadOffset: int64(1000),
		AssetsDir:      "assets/static",
	}

	cfg.IgnoreSuffix = append(cfg.IgnoreSuffix, ".jar", ".war") // code
	cfg.IgnoreSuffix = append(cfg.IgnoreSuffix, ".html", ".js", ".css", ".java", ".class")
	cfg.IgnoreSuffix = append(cfg.IgnoreSuffix, ".gz", ".tar", ".zip", ".rar")
	cfg.IgnoreSuffix = append(cfg.IgnoreSuffix, ".jpg", ".png", ".xls", ".xlxs", ".pdf")

	return cfg

}

func UnmarshalServerConfFromIni(defaultCfg *ServerCommonConf, content string) (cfg *ServerCommonConf, err error) {
	cfg = defaultCfg
	if cfg == nil {
		cfg = GetDefaultServerConf()
	}

	conf, err := ini.Load(strings.NewReader(content))
	if err != nil {
		err = fmt.Errorf("parse ini conf file error: %v", err)
		return nil, err
	}

	var (
		tmpStr string
		ok     bool
		v      int64
	)
	if tmpStr, ok = conf.Get("common", "bind_addr"); ok {
		cfg.BindAddr = tmpStr
	}

	if tmpStr, ok = conf.Get("common", "bind_port"); ok {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err != nil {
			err = fmt.Errorf("Parse conf error: invalid bind_port")
			return
		} else {
			cfg.BindPort = int(v)
		}
	}

	if tmpStr, ok = conf.Get("common", "conn_max_time"); ok {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err != nil {
			err = fmt.Errorf("Parse conf error: invalid con_max_time")
			return
		} else {
			cfg.ConnMaxTime = time.Duration(v) * time.Minute
		}
	}

	if tmpStr, ok = conf.Get("common", "heart_interval"); ok {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err != nil {
			err = fmt.Errorf("Parse conf error: invalid heart_interval")
			return
		} else {
			cfg.HeartInterval = time.Duration(v) * time.Second
		}
	}

	if tmpStr, ok = conf.Get("common", "log_file"); ok {
		cfg.LogFile = tmpStr
		if cfg.LogFile == "console" {
			cfg.LogWay = "console"
		} else {
			cfg.LogWay = "file"
		}
	}

	if tmpStr, ok = conf.Get("common", "log_level"); ok {
		cfg.LogLevel = tmpStr
	}

	if tmpStr, ok = conf.Get("common", "log_max_days"); ok {
		v, err = strconv.ParseInt(tmpStr, 10, 64)
		if err == nil {
			cfg.LogMaxDays = v
		}
	}

	if tmpStr, ok = conf.Get("common", "user"); ok {
		cfg.User = tmpStr
	}

	if tmpStr, ok = conf.Get("common", "pwd"); ok {
		cfg.Pwd = tmpStr
	}

	if tmpStr, ok = conf.Get("common", "ignore_suffix"); ok {
		cfg.IgnoreSuffix = ParseIgnoreSuffix(tmpStr)
	}

	if tmpStr, ok = conf.Get("common", "last_read_offset"); ok {
		v, err = strconv.ParseInt(tmpStr, 10, 64)
		if err != nil {
			err = fmt.Errorf("Parse conf error: invalid last read offset")
			return
		} else {
			cfg.LastReadOffset = v
		}
	}

	if tmpStr, ok = conf.Get("common", "assets_dir"); ok {
		cfg.AssetsDir = tmpStr
	}
	return cfg, nil
}

func ParseIgnoreSuffix(ignore string) []string {
	if ignore == "" {
		return []string{}
	} else {
		return strings.Split(ignore, ",")
	}
}

func init() {
	glbEnvs = make(map[string]string)
	envs := os.Environ()
	for _, env := range envs {
		kv := strings.Split(env, "")
		if len(kv) == 2 {
			glbEnvs[kv[0]] = kv[1]
		}
	}
}

type Values struct {
	Envs map[string]string
}

func GetValues() *Values {
	return &Values{
		Envs: glbEnvs,
	}
}

func ReaderContent(in string) (out string, err error) {
	tmpl, errRet := template.New("log-tail").Parse(in)
	if errRet != nil {
		err = errRet
		return
	}
	buffer := bytes.NewBufferString("")
	v := GetValues()
	err = tmpl.Execute(buffer, v)
	if err != nil {
		return
	}
	out = buffer.String()
	return
}

func GetRenderedConfFromFile(path string) (out string, err error) {
	var b []byte
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	content := string(b)
	out, err = ReaderContent(content)
	return
}
