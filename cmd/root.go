package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log-tail/g"
	"log-tail/models/config"
	"log-tail/server"
	"log-tail/util/log"
	"log-tail/util/version"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfgFile       string
	showVersion   bool
	bindAddr      string
	bindPort      int
	conMaxTime    int64
	heartInterval int64
	logFile       string
	logLevel      string
	logMaxDays    int64
	user          string
	pwd           string
	ignoreSuffix  string
)

const (
	CfgFileTypeIni = iota
	CfgFileTypeCmd
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "ci", "c", "./log_tail.ini", "ci file of log-tail")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of log-tail")
}

var rootCmd = &cobra.Command{
	Use:   "log-tail",
	Short: "log tail in web",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		content, err := config.GetRenderedConfFromFile(cfgFile)
		if err != nil {
			return err
		}
		g.GlbServerCfg.CfgFile = cfgFile
		err = parseServerCommonCfg(CfgFileTypeIni, content)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = runServer()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func parseServerCommonCfg(fileType int, content string) (err error) {
	if fileType == CfgFileTypeIni {
		err = parseServerCommonCfgFromIni(content)
	} else if fileType == CfgFileTypeCmd {
		err = parseServerCommonCfgFromCmd()
	}

	if err != nil {
		return
	}
	return nil
}

func parseServerCommonCfgFromIni(content string) (err error) {
	cfg, err := config.UnmarshalServerConfFromIni(&g.GlbServerCfg.ServerCommonConf, content)
	if err != nil {
		return err
	}
	g.GlbServerCfg.ServerCommonConf = *cfg
	return
}

func parseServerCommonCfgFromCmd() (err error) {
	g.GlbServerCfg.BindAddr = bindAddr
	g.GlbServerCfg.BindPort = bindPort
	g.GlbServerCfg.ConnMaxTime = time.Duration(conMaxTime) * time.Minute
	g.GlbServerCfg.HeartInterval = time.Duration(heartInterval) * time.Second
	g.GlbServerCfg.LogFile = logFile
	g.GlbServerCfg.LogLevel = logLevel
	g.GlbServerCfg.LogMaxDays = logMaxDays
	g.GlbServerCfg.User = user
	g.GlbServerCfg.Pwd = pwd

	if ignoreSuffix != "" {
		g.GlbServerCfg.IgnoreSuffix = config.ParseIgnoreSuffix(ignoreSuffix)
	}

	if logFile == "console" {
		g.GlbServerCfg.LogWay = "console"
	} else {
		g.GlbServerCfg.LogWay = "file"
	}
	return
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServer() (err error) {
	log.InitLog(g.GlbServerCfg.LogWay, g.GlbServerCfg.LogFile, g.GlbServerCfg.LogLevel, g.GlbServerCfg.LogMaxDays)

	svr, err := server.NewService()
	if err != nil {
		return err
	}
	svr.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	return nil
}
