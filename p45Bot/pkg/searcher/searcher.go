package searcher

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SearchFiles(root, ext string, excl []string, v bool) []string {
	dir := os.DirFS(root)
	//the below recurses from the root to find all file with teh matching filetype extension form the manifest file
	all, err := findFiles(dir, root, ext)
	if err != nil {
		log.Println("Error walking directory -", err)
		return nil
	}
	var filtered []string
	//from the that list we then filter out the exclusions
	for _, a := range all {
		valid := filterFiles(a, ext, excl, v)
		if valid == 0 {
			filtered = append(filtered, a)
		}
	}
	return filtered
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

func filterFiles(f, ext string, excl []string, v bool) int {
	//return value and where we keep track of hits against the exclusions
	i := 0
	//loop over the exclusions we have form the manifest file
	for _, e := range excl {
		//munging the exclusion with the file type extension
		ex := e + ext
		//stripping the full file path back to just the last piece - i.e. the filename with extension so we can compare
		fName := filepath.Base(f)
		if strings.EqualFold(fName, ex) {
			if v {
				fmt.Println("The strings match")
			}
			i++
		} else {
			if v {
				fmt.Println("Strings do not match")
			}
		}
	}
	return i
}
