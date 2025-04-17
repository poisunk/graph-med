package conf

import (
	"gopkg.in/yaml.v2"
	"os"
)

const (
	DefaultMedicalAttrPath = "configs/medical_attr.yaml"
)

type MedicalAttr struct {
	MedicalAttrs map[string]string `yaml:"medical_attrs"`
}

func LoadMedicalAttrs(path string) (map[string]string, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var attrs MedicalAttr
	err = yaml.Unmarshal(yamlFile, &attrs)
	return attrs.MedicalAttrs, err
}
