/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"calfinn.io/p45bot/pkg/opts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Verbose bool
var DryRun bool
var Manifest string
var Scanpath string
var Output string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "p45bot",
	Short: "P45Bot scans a given path for matching files, strings and then validates users exist.",
	Long: `P45Bot scans a given path for matching files, strings and then validates users exist.
It is designed with Terraform files in mind and currently only geared towards those and using 
AzureAD as a source of truth but future iterations will expand this.

Viper is present to handle the configuration loading and operates against a local
.p45bot config file containg paths to scan, where the manifest is and what the 
credentials are for connecting to AzureAD.

The manifest file requires a minimum config of the fileType, search strings and 
filename exclusions (e.g. variables).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//PersistentPreRun: func(cmd *cobra.Command, args []string) {
	//	// Bind the persistent flag to the Viper value
	//	viper.BindPFlag("manifest", cmd.PersistentFlags().Lookup("manifest"))
	//},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Manifest:", Manifest)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.p45bot.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dryrun", "d", false, "dryrun only")
	rootCmd.PersistentFlags().StringVarP(&Manifest, "manifest", "m", "./manifests/example.json", "Path to manifest file")
	rootCmd.PersistentFlags().StringVarP(&Scanpath, "scanpath", "s", "./", "Path to scan")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "", "Output type for commands")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("manifest", rootCmd.PersistentFlags().Lookup("manifest"))
	viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))
	viper.BindPFlag("scanpath", rootCmd.PersistentFlags().Lookup("scanpath"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".p45bot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".p45bot")
	}

	viper.AutomaticEnv() // read in environment variables that
	viper.BindEnv("directoryconfig.clientsecret", "AZ_SP_SECRET")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
