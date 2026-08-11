package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Version  string    `json:"version"`
	Packages []Package `json:"packages"`
	SQL      []SQL     `json:"sql"`
}

type Package struct {
	Name                      string   `json:"name"`
	Engine                    string   `json:"engine"`
	Schema                    Paths    `json:"schema"`
	Queries                   Paths    `json:"queries"`
	Gen                       Gen      `json:"gen"`
	Codegen                   []Codegen `json:"codegen"`
	StrictFunctionSignatures  bool     `json:"strict_function_signatures"`
}

type Codegen struct {
	Out     string `json:"out"`
	Plugin  string `json:"plugin"`
	Options []byte `json:"options"`
}

type Gen struct {
	Go   *Go   `json:"go,omitempty"`
	JSON *JSON `json:"json,omitempty"`
}

type Go struct {
	Package                    string   `json:"package"`
	Out                        string   `json:"out"`
	EmitInterface              bool     `json:"emit_interface"`
	EmitJSONTags               bool     `json:"emit_json_tags"`
	EmitDBTags                 bool     `json:"emit_db_tags"`
	EmitPreparedQueries        bool     `json:"emit_prepared_queries"`
	EmitExactTableNames        bool     `json:"emit_exact_table_names"`
	EmitEmptySlices            bool     `json:"emit_empty_slices"`
	EmitExportedQueries        bool     `json:"emit_exported_queries"`
	EmitMethodsWithDBArgument  bool     `json:"emit_methods_with_db_argument"`
	EmitPointersForNullTypes   bool     `json:"emit_pointers_for_null_types"`
	EmitEnumValidMethod        bool     `json:"emit_enum_valid_method"`
	EmitAllEnumValues          bool     `json:"emit_all_enum_values"`
	JSONTagsCaseStyle          string   `json:"json_tags_case_style"`
	SQLPackage                 string   `json:"sql_package"`
	SQLDriver                  string   `json:"sql_driver"`
	Overrides                  []Override `json:"overrides"`
	Rename                     map[string]string `json:"rename"`
	OutputBatchFileName        string   `json:"output_batch_file_name"`
	OutputDBFileName           string   `json:"output_db_file_name"`
	OutputModelsFileName       string   `json:"output_models_file_name"`
	OutputQuerierFileName      string   `json:"output_querier_file_name"`
	OutputCopyFromFileName     string   `json:"output_copyfrom_file_name"`
}

type JSON struct {
	Out string `json:"out"`
}

type Override struct {
	DBType   string `json:"db_type"`
	GoType   GoType `json:"go_type"`
	Nullable bool   `json:"nullable"`
	Column   string `json:"column"`
	Table    string `json:"table"`
}

type GoType struct {
	Import string `json:"import"`
	Type   string `json:"type"`
}

type SQL struct {
	Engine  string   `json:"engine"`
	Schema  Paths    `json:"schema"`
	Queries Paths    `json:"queries"`
	Gen     SQLGen   `json:"gen"`
	Codegen []Codegen `json:"codegen"`
	StrictFunctionSignatures bool `json:"strict_function_signatures"`
}

type SQLGen struct {
	Go   *Go   `json:"go,omitempty"`
	JSON *JSON `json:"json,omitempty"`
}

type Paths []string

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
				return fileI < fileJ
			}
		}

		if len(matchI) > 0 && len(matchJ) == 0 {
			return true
		}
		if len(matchI) == 0 && len(matchJ) > 0 {
			return false
		}

		return fileI < fileJ
	})
}

func (c *Config) ResolveSchema() error {
	for idx, s := range c.SQL {
		var files []string
		for _, path := range s.Schema {
			info, err := os.Stat(path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil
					}
					if strings.HasSuffix(p, ".sql") {
						files = append(files, p)
					}
					return nil
				})
				if err != nil {
					return err
				}
			} else {
				files = append(files, path)
			}
		}
		sortFiles(files)
		c.SQL[idx].Schema = files
	}
	return nil
}
