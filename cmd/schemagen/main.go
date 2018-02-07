package main

import (
	"context"
	"fmt"
	"os"

	"github.com/burdiyan/schemagen"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func main() {
	var (
		cfgFile string
		cfg     schemagen.Config
	)

	rootCmd := &cobra.Command{
		Use: "schemagen",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cfgFile, &cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return schemagen.Run(context.TODO(), cfg)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default .schemagen.yaml in the current directory")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func initConfig(cfgFile string, cfg *schemagen.Config) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".schemagen")
	}

	if err := viper.ReadInConfig(); err != nil {

	}

	viper.Unmarshal(cfg)

	return nil
}
