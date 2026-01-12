package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile string
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
)

var rootCmd = &cobra.Command{
	Use:   "ufi003-cli",
	Short: "A brief description of your application",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
		initConfig()
	},
}

func initLogger() {
	var err error

	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	sugar = logger.Sugar()
}

func initConfig() {
	if cfgFile == "" {
		sugar.Fatal("config file must be specified using --config flag")
	}

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		sugar.Fatalf("failed to read config: %v", err)
	}

	sugar.Infow("configuration loaded",
		"config_file", viper.ConfigFileUsed(),
	)
}

func Execute() {
	defer func() {
		_ = logger.Sync()
	}()

	err := rootCmd.Execute()
	if err != nil {
		sugar.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ufi003-cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
