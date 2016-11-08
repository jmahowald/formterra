package tfproject

import (
	log "github.com/Sirupsen/logrus"

	"github.com/ghodss/yaml"
)

func (t *TerraformProjectSkeleton) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Warn("err: %v\n", err)
		return err
	}
	return nil
}

func (t *TerraformProjectSkeleton) MarshalYAML() ([]byte, error) {
	yaml, err := yaml.Marshal(t)
	return yaml, err
}
