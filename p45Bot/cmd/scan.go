/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"calfinn.io/p45bot/pkg/azure"
	"calfinn.io/p45bot/pkg/opts"
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
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		//var results []searcher.SearchForResults
		r := searcher.SearchForFiles(
			viper.GetString("scanpath"),
			y.GetString("fileType"),
			y.GetStringSlice("fileNameExclusions"))

		if err != nil {
			panic(fmt.Errorf("fatal error string search file: %w", err))
		}
		client, err := azure.NewAZClient(
			viper.GetString("directoryconfig.tenantid"),
			viper.GetString("directoryconfig.clientid"),
			viper.GetString("directoryconfig.clientsecret"))
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		rs, err := searcher.SearchString(
			viper.GetString("scanpath"),
			r.MatchingFiles,
			y.GetStringSlice("searchStrings"))
		if err != nil {
			fmt.Println(err)
		}
		uniqueList := searcher.UniqueUpns(rs)
		for i := range uniqueList {
			userSearch, err := azure.GetUserByUPN(client, uniqueList[i].Upn)
			if err != nil {
				fmt.Println(err)
			}
			exists, err := azure.CheckUserExists(userSearch)
			if err != nil {
				fmt.Println(err)
			}
			uniqueList[i].Exists = exists
			if opts.GetVerbose() {
				fmt.Printf("%+v\n", uniqueList[i])
			}
		}

		switch Output {
		case "":
			fmt.Println("File search stats:")
			searcher.PrettyPrintJson(r)
			fmt.Println("Raw Results (no UPN checks):")
			searcher.PrettyPrintJson(rs)
			fmt.Println("Unique UPNs, validated and collated:")
			searcher.PrettyPrintJson(uniqueList)
			//fmt.Printf("%+v\n", unqiueList)
			//fmt.Printf("%+v\n", rs)

		case "json":
			data := searcher.DataOutputs{
				Stats:        *r,
				Raw:          rs,
				ValidatedUpn: uniqueList,
			}

			file, _ := json.MarshalIndent(data, "", " ")
			_ = ioutil.WriteFile("data.json", file, 0644)
			//searcher.OutputToJson(r)
			//searcher.OutputToJson(rs)
		default:
			searcher.PrettyPrintJson(r)
			//	searcher.PrettyPrintJson(rs)
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
