package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "ufi003-cli",
	Short: "A brief description of your application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ufi003-cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	// 查找配置文件
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 当前目录
		viper.AddConfigPath(".")

		// 用户家目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)

		viper.SetConfigName(".ufi003-cli.yaml")
		viper.SetConfigType("yaml")
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}
}
