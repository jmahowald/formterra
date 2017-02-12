package tfproject

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

// UnmarshalYAML populate a skeleton from yaml
func (t *TerraformProjectSkeleton) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Warn("err: %v\n", err)
		return err
	}
	return nil
}

// MarshalYAML Write skeleton out to yaml
func (t *TerraformProjectSkeleton) MarshalYAML() ([]byte, error) {
	yaml, err := yaml.Marshal(t)
	return yaml, err
}

// MarshalYAML Write skeleton out to yaml
func (t *TerraformModuleDefinition) MarshalYAML() ([]byte, error) {
	yaml, err := yaml.Marshal(t)
	return yaml, err
}

// UnmarshalYAML parses out a project definition
func (t *TerraformModuleDefinition) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Warn("err: %v\n", err)
		return err
	}
	return nil
}

// Utility function to make testing slightly easier
func skelFromYaml(data []byte) (*TerraformProjectSkeleton, error) {
	t := TerraformProjectSkeleton{}
	err := t.UnmarshalYAML(data)
	return &t, err
}
