package main

import (
	"fmt"
	"github.com/arugal/log-tail/g"
	"github.com/arugal/log-tail/server"
	"github.com/arugal/log-tail/util/log"
	"github.com/arugal/log-tail/util/version"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile     string
	showVersion bool
	bindAddr    string
	bindPort    int
	logLevel    string
)

const (
	originalCmd = "-"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file of log-tail")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of log-tail")
	rootCmd.PersistentFlags().StringVarP(&bindAddr, "host", "H", "-", "host of log-tail")
	rootCmd.PersistentFlags().IntVarP(&bindPort, "port", "p", -1, "port of log-tail")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "-", "log level of log-tail")
}

var rootCmd = &cobra.Command{
	Use:   "log-tail",
	Short: "log tail in web",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		err := g.Load(cfgFile)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cnfOverrideFromCmd()

		err = runServer()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func cnfOverrideFromCmd() {
	if bindAddr != originalCmd {
		g.ServerCnf.Host = bindAddr
	}

	if bindPort > 0 {
		g.ServerCnf.Port = bindPort
	}

	if logLevel != originalCmd {
		g.CommonCnf.Log.Level = logLevel
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServer() (err error) {
	logCnf := g.CommonCnf.Log
	log.InitLog(logCnf.Way, logCnf.File, logCnf.Level, logCnf.MaxDays)

	svr, err := server.NewService()
	if err != nil {
		return err
	}
	svr.Start()
	return
}
