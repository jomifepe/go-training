package cmd

import (
	"fmt"
	"github.com/jomifepe/gin_api/api"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gin_api",
		Short: "Gin_api is a test api that uses the gin framework",
		Long: `Gin_api is a test api developed for training. It uses the gin framework to handle the requests, 
jwt tokens for the authentication and viper & cobra for cli commands and configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			logging.NewLogger()
			api.Start(viper.GetString("API_PORT"))
		},
	}
	cfgFile string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("LOG_LEVEL", "error")
	viper.SetDefault("LOG_FORMAT_JSON", false)
	viper.SetDefault("API_PORT", "3000")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_DB", "go_test")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("JWT_ACCESS_SECRET", util.GetRandStringBytes(64))

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to an environment config file")
}

func initConfig() {
	if len(cfgFile) > 0 /* a config file path was passed in */ {
		fmt.Println(cfgFile)
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			logrus.New().WithFields(logrus.Fields{
				"error": err.Error(),
				"file": cfgFile,
			}).Errorln("[CFG] Failed to read config file")
		}
		return
	}

	if util.FileExists(".env") /* config file on the root path */ {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			logrus.New().WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorln("[CFG] Failed to read config file")
		}
		return
	}

	/* no config file or flag */
	content := "# Created on the first run with default values\n\n"
	for key, value := range viper.AllSettings() {
		content += fmt.Sprintf("%v=%v\n", strings.ToUpper(key), value)
	}
	if err := util.AppendToFile(".env", content); err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"file": ".env",
			"error": err.Error(),
		}).Errorln("[IO] Failed to write string to file")
	} else {
		if abs, err := filepath.Abs(".env"); err != nil {
			logrus.New().WithFields(logrus.Fields{
				"file": abs,
			}).Infoln("[IO] No config input or file found, writing defaults...")
		}
	}
}