package searcher

import (
	"log"
	"reflect"
	"sort"
	"testing"
	"testing/fstest"
)

func TestFindFiles(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		ext      string
		files    fstest.MapFS
		expected []string
	}{
		{
			name: "no files with extension",
			root: ".",
			ext:  ".txt",
			files: fstest.MapFS{
				"file.go":             {},
				"afile1.go":           {},
				"dir1/file2.py":       {},
				"dir1/dir2/file4.php": {},
				"subfolder2/file.go":  {},
			},
			expected: nil,
		},
		{
			name: "single file with extension",
			root: ".",
			ext:  ".txt",
			files: fstest.MapFS{
				"file.go":             {},
				"afile1.go":           {},
				"dir1/file2.txt":      {},
				"dir1/dir2/file4.php": {},
				"subfolder2/file.go":  {},
			},
			expected: []string{"dir1/file2.txt"},
		},
		{
			name: "multiple files with extension",
			root: ".",
			ext:  ".txt",
			files: fstest.MapFS{
				"file.go":                   {},
				"afile1.go":                 {},
				"dir1/file2.txt":            {},
				"dir1/dir2/file4.txt":       {},
				"subfolder2/file.php":       {},
				"subfolder3/dir4/file5.txt": {},
			},
			expected: []string{"dir1/file2.txt", "dir1/dir2/file4.txt", "subfolder3/dir4/file5.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := fstest.MapFS(tt.files)

			// call the function being tested
			got, err := findFiles(fsys, tt.root, tt.ext)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			sort.Strings(got)
			sort.Strings(tt.expected)
			log.Printf(" files found: got %v, want %v", got, tt.expected)
			// check that the correct files were found
			if reflect.DeepEqual(got, tt.expected) != true {
				t.Errorf("unexpected files found: got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilterFiles(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		ext      string
		excl     []string
		verbose  bool
		expected int
	}{
		{
			name:     "No Exclusions",
			filePath: "subfolder3/dir4/file5.txt",
			ext:      ".txt",
			excl:     []string{"terraform-configuration"},
			verbose:  false,
			expected: 0,
		},
		{
			name:     "Exclusions",
			filePath: "subfolder3/dir4/terraform-configuration.tf",
			ext:      ".tf",
			excl:     []string{"terraform-configuration"},
			verbose:  false,
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterFiles(tt.filePath, tt.ext, tt.excl, tt.verbose)
			if got != tt.expected {
				t.Errorf("unexpected exclusion found: got %v, want %v", got, tt.expected)
			}
		})
	}
}
