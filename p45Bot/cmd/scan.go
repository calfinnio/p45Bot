/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"calfinn.io/p45bot/pkg/searcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scan called")
		y := viper.New()
		y.SetConfigFile(viper.GetString("manifest"))

		if viper.GetBool("verbose") {
			fmt.Println("Using manifest file at:", viper.GetString("manifest"))
		}
		err := y.ReadInConfig() // Find and read the config file
		if err != nil {         // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		for _, s := range searcher.SearchFiles(viper.GetString("scanpath"), y.GetString("fileType"), y.GetStringSlice("fileNameExclusions"), viper.GetBool("verbose")) {
			println(s)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
