package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/fatih/structs"
	cliYTY "github.com/jkandasa/jenkinsctl/pkg/types/cli"
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

func GetResource(bytes []byte) (interface{}, error) {
	kindData := &cliYTY.Kind{}
	err := yaml.Unmarshal(bytes, kindData)
	if err != nil {
		return nil, err
	}

	var resource interface{}

	switch kindData.Kind {
	case cliYTY.KindTypeBuild:
		resource = &cliYTY.KindBuild{}

	case cliYTY.KindTypeJob:
		resource = &cliYTY.KindJob{}

	default:
		return nil, fmt.Errorf("unknown kind:%s", kindData.Kind)
	}

	if err = yaml.Unmarshal(bytes, resource); err != nil {
		return nil, err
	}
	return resource, nil
}
