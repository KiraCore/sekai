package functions

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FunctionMetaField is a struct for each parameter in listing functions
type FunctionMetaField struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// FunctionMeta is a struct for each functions in listing functions
type FunctionMeta struct {
	Description string                       `json:"description"`
	Parameters  map[string]FunctionMetaField `json:"parameters"`
}

// GetAllFunctions is a function to get all functions registered
func GetAllFunctions() map[string]FunctionMeta {
	functions := map[string]FunctionMeta{}

	err := filepath.Walk("functions",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".json") {
				file, _ := ioutil.ReadFile(path)
				functionsFromFile := map[string]FunctionMeta{}

				err := json.Unmarshal([]byte(file), &functionsFromFile)
				if err == nil {
					for k, v := range functionsFromFile {
						functions[k] = v
					}
				}
			}

			return nil
		})

	if err != nil {
		return functions
	}

	return functions
}

// AllFunctions is all functions registered
var AllFunctions = GetAllFunctions()
