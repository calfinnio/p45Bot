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
	Short: "Scans file extensions at path for strings and then validates those exist in AD.",
	Long: `With a correctly defined manfiest, configuration and environment the following command
	should be all that is required:
	
	    p45bot scan
	
	This will default to outputting the results to the console.  As an alternative a concatenated 
	collection of results can be dumped to json using the --output flag:

	    p45bot scan --output json
	
	The current flow is:
	- scan path for all files matching provided extension
	- filter out those that are defined in the manifest exclusions (as an example you might want
	to exclude all variables.tf files from the scrape)
	- Iterates through those files and scans for matching strings/regex
	- Sorts those results and structures them by unique UPNs
	- Connects to AzureAD and searches for the UPN.  Returns true/false for each
	`,
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
