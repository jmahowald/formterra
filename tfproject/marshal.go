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
