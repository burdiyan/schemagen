package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/burdiyan/schemagen"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	rootCmd.PersistentFlags().BoolVar(&cfg.NoFetch, "no-fetch", false, "disable fetch and only compile schemas from source directory")

	if err := bindFlags(rootCmd.PersistentFlags(), viper.GetViper()); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func bindFlags(fs *pflag.FlagSet, v *viper.Viper) error {
	var err error

	fs.VisitAll(func(f *pflag.Flag) {
		err = v.BindPFlag(strings.Replace(f.Name, "-", "_", -1), f)
		if err != nil {
			return
		}
	})

	return err
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
