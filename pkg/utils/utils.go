package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/fatih/structs"
	cliML "github.com/jkandasa/jenkinsctl/pkg/model/cli"
	"gopkg.in/yaml.v2"
)

// StructToMap converts struct to a map
func StructToMap(data interface{}) map[string]interface{} {
	return structs.Map(data)
}

// FileToStruct loads the file date to the given struct
func FileToStruct(filename string, out interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, out)
}

func GetResource(filename string) (interface{}, error) {
	kindData := &cliML.Kind{}
	err := FileToStruct(filename, kindData)
	if err != nil {
		return nil, err
	}

	var resource interface{}

	switch kindData.Kind {
	case cliML.KindTypeBuild:
		resource = &cliML.KindBuild{}

	case cliML.KindTypeJob:
		resource = &cliML.KindJob{}

	default:
		return nil, fmt.Errorf("unknown kind:%s", kindData.Kind)
	}

	if err = FileToStruct(filename, resource); err != nil {
		return nil, err
	}
	return resource, nil
}
