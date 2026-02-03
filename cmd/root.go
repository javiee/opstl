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
	Short: "CLI tool to automate Kubernetes operations",
	Long: `opstl - Operations Tool

A CLI tool to automate Kubernetes tasks for your team's namespaces.

Commands:
  kubernetes (k8, kubectl)    Kubernetes operations
    get-pods                  List all pods owned by the configured team
    get-events                List all warning events in owned namespaces
    watch-pods                Watch and stream logs from pods in real-time

Examples:
  opstl kubernetes get-pods
  opstl k8 get-pods --status Running
  opstl k8 get-events
  opstl k8 watch-pods --label-selector app=myapp --regex "myapp-.*"

Configuration:
  Place a utils.yaml file in the current directory, $HOME/.utils, or /etc/
  with your team's namespace configuration.`,
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
}
