package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "opstl",
	Short: "This cli is thought to automate multiple tasks",
	Long: `
	Usage:
	kubectl `,
	// The main action executed when no subcommand is provided:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from root command!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you can define persistent/global flags, for example:
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opstl.yaml)")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("utils")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$HOME/.utils")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Print(viper.Get("filteredNamespaces"))
}
