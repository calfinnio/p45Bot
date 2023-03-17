package searcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"calfinn.io/p45bot/pkg/opts"
)

type SearchForResults struct {
	Root               string   `json:"root"`
	FileType           string   `json:"filetype"`
	TotalFilesCount    int      `json:"totalfilescount"`
	MatchingFilesCount int      `json:"matchingfilescount"`
	MatchingFiles      []string `json:"matchingfiles"`
}

type SearchStringResult struct {
	FileName     string
	SearchString string
	Upn          string
	LineNumber   int
	Exists       bool
}

type SearchStringResults []SearchStringResult

type searchString struct {
	SearchString string
	FileHit      FileHit
}

type FileHit struct {
	FileName   string
	LineNumber int
}
type UpnFilteredResult struct {
	Upn           string
	SearchStrings []searchString
	Exists        bool
}

type UpnFilteredResults []UpnFilteredResult

type CliOutput interface {
	prettyPrint() string
	outputJson()
}
type DataOutputs struct {
	Stats        SearchForResults
	Raw          SearchStringResults
	ValidatedUpn UpnFilteredResults
}

//func PrettyPrintJson() {
//	prettyJSON, err := json.MarshalIndent(r, "", "  ")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//}

func PrettyPrintJson(c CliOutput) {
	fmt.Println(c.prettyPrint())
}

func OutputToJson(c CliOutput) {
	c.outputJson()
}

func SearchForFiles(root, ext string, excl []string) *SearchForResults {
	results := &SearchForResults{}
	results.Root = root
	results.FileType = ext

	dir := os.DirFS(root)
	//the below recurses from the root to find all file with teh matching filetype extension form the manifest file
	all, err := findFiles(dir, root, ext)
	if err != nil {
		log.Println("Error walking directory -", err)
		return nil
	}
	//Adding count
	results.TotalFilesCount = len(all)
	var filtered []string
	//from the that list we then filter out the exclusions
	for _, a := range all {
		valid := filterFiles(a, ext, excl)
		if valid == 0 {
			filtered = append(filtered, a)
		}
	}
	results.MatchingFilesCount = len(filtered)
	results.MatchingFiles = filtered
	return results
}

func (s SearchForResults) prettyPrint() string {
	// VERY custom logic for generating an english greeting
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	//fmt.Println(string(prettyJSON))
	return string(prettyJSON)
}

func (s SearchForResults) outputJson() {
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(s); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON file created successfully")
}

func SearchString(root string, files, search []string) (SearchStringResults, error) {
	r := SearchStringResults{}
	//openBracket := "["
	//closeBracket := "]"
	for _, file := range files {
		targetFile := filepath.Join(root, file)
		fmt.Println("Processing", targetFile)
		openFile, err := os.Open(targetFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer openFile.Close()
		patterns := make(map[string]*regexp.Regexp)
		for _, target := range search {
			pattern := regexp.MustCompile(regexp.QuoteMeta(target) + `\s*=\s*\[(.*?)\]`)
			//fmt.Println(pattern)
			patterns[target] = pattern
		}

		scanner := bufio.NewScanner(openFile)
		lineNumber := 1
		replacer := strings.NewReplacer(`"`, "", ` `, "")
		for scanner.Scan() {
			line := scanner.Text()
			// Check each pattern and find matches
			for target, pattern := range patterns {
				//fmt.Println("Running regex")
				matches := pattern.FindStringSubmatch(line)
				//fmt.Println("Number of matches:", len(matches))
				//fmt.Println(matches)
				if len(matches) > 1 {
					textBetweenBrackets := matches[1]
					for _, f := range strings.Split(textBetweenBrackets, ",") {
						var t SearchStringResult
						t.FileName = targetFile
						t.SearchString = string(target)
						t.Upn = replacer.Replace(f)
						t.LineNumber = lineNumber
						r = append(r, t)
					}
					if opts.GetVerbose() {
						fmt.Printf("For '%s', text between brackets: %s\n", target, textBetweenBrackets)
					}

				}
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
		//fmt.Println(found)
		//split result in to upns
		//populate struct
		//concat structs
	}
	return r, nil
}

func (s SearchStringResults) prettyPrint() string {
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(prettyJSON)
}

func (s SearchStringResults) outputJson() {
	file, err := os.Create("results.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(s); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON file created successfully")
}
func UniqueUpns(s SearchStringResults) UpnFilteredResults {
	e := getUniqueUpn(s)
	if opts.GetVerbose() {
		fmt.Println(e)
	}
	var UpnsFiltered UpnFilteredResults
	for _, u := range e {
		f := filterSearchResultByUPN(s, u)
		UpnsFiltered = append(UpnsFiltered, f)
	}
	return UpnsFiltered
}

func (s UpnFilteredResults) prettyPrint() string {
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(prettyJSON)
}

func (s UpnFilteredResults) outputJson() {
	file, err := os.Create("results.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(s); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON file created successfully")
}

func findFiles(dir fs.FS, root, ext string) ([]string, error) {

	var a []string
	err := fs.WalkDir(dir, ".", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func filterFiles(f, ext string, excl []string) int {
	//return value and where we keep track of hits against the exclusions
	i := 0
	//loop over the exclusions we have form the manifest file
	for _, e := range excl {
		//munging the exclusion with the file type extension
		ex := e + ext
		//stripping the full file path back to just the last piece - i.e. the filename with extension so we can compare
		fName := filepath.Base(f)
		if strings.EqualFold(fName, ex) {
			if opts.GetVerbose() {
				fmt.Printf("Filename is excluded - %v (input) and %v (excluded)\n", fName, ex)
			}
			i++
		} else {
			if opts.GetVerbose() {
				fmt.Printf("Filename not in exclusion list - %v (input) and %v (excluded) \n", fName, ex)
			}
		}
	}
	return i
}

func getUniqueUpn(s SearchStringResults) []string {
	uniqueUpns := map[string]bool{}

	for _, p := range s {
		uniqueUpns[p.Upn] = true
	}

	upns := []string{}
	for n := range uniqueUpns {
		upns = append(upns, n)
	}

	return upns
}

func filterSearchResultByUPN(s SearchStringResults, upn string) UpnFilteredResult {
	f := UpnFilteredResult{}
	f.Upn = upn
	//sPtr := &s
	var ss []searchString
	for _, p := range s {
		if p.Upn == upn {
			fh := &FileHit{
				FileName:   p.FileName,
				LineNumber: p.LineNumber,
			}
			h := &searchString{
				SearchString: p.SearchString,
				FileHit:      *fh,
			}
			ss = append(ss, *h)
		}
		f.SearchStrings = ss
	}
	if opts.GetVerbose() {
		fmt.Println(f)
	}
	return f
}
