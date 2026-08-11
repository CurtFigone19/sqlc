package migrations

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var numericPrefixRegex = regexp.MustCompile(`^(\d+)`)

func sortFiles(files []string) {
	sort.Slice(files, func(i, j int) bool {
		fileI := filepath.Base(files[i])
		fileJ := filepath.Base(files[j])

		matchI := numericPrefixRegex.FindStringSubmatch(fileI)
		matchJ := numericPrefixRegex.FindStringSubmatch(fileJ)

		if len(matchI) > 0 && len(matchJ) > 0 {
			valI, errI := strconv.ParseUint(matchI[1], 10, 64)
			valJ, errJ := strconv.ParseUint(matchJ[1], 10, 64)
			if errI == nil && errJ == nil {
				if valI != valJ {
					return valI < valJ
				}
				// If numeric values are equal (e.g. 01_init.sql and 1_init.sql),
				// fallback to lexicographical sort of the whole filename to be consistent.
				return fileI < fileJ
			}
		}

		// Fallback to lexicographical sorting if one or both do not have numeric prefixes
		if len(matchI) > 0 && len(matchJ) == 0 {
			return true // numeric prefixes come first
		}
		if len(matchI) == 0 && len(matchJ) > 0 {
			return false
		}

		return fileI < fileJ
	})
}

func Read(fsys fs.FS, dir string) ([]string, error) {
	var files []string
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortFiles(files)
	return files, nil
}
