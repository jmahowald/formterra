package tfproject

import (
	"os"
	"path/filepath"

	"fmt"

	"github.com/backpack/formterra/core"
	//TODO wrap this up in our own to allow us to switch out easier?
	log "github.com/Sirupsen/logrus"
)

const dirMode = 0755

//TerraformProjectRequest Can create a terraform layer for u
type TerraformProjectRequest interface {
	Create() (TerraformLayer, bool)
}

// TfConfig configuration for this project.  By default uses viper
var TfConfig = viperConfig

func (t TerraformLayer) getDir() (string, bool) {
	fullPath, exists := layerExists(t.Name)
	log.Debug("Attempting to create:", fullPath)
	if !exists {
		err := os.MkdirAll(fullPath, dirMode)
		if err != nil {
			log.Fatalf("Unable to create directory at %s:%s", fullPath, err)
		}
	}
	return fullPath, exists
}

func formTerraVersion() string {
	return fmt.Sprintf("Version: %s  BuildTime:%s", core.Version, core.BuildTime)
}

func layerExists(path string) (string, bool) {
	dirPath := filepath.Join(getString(TerraformDir),
		getString("env"), path)
	_, err := os.Stat(dirPath)
	return dirPath, (err == nil)
}

func (t *TerraformLayer) dir() (string, bool) {
	dirPath := filepath.Join(getString(TerraformDir),
		getString("env"), t.Name)

	_, err := os.Stat(dirPath)
	return dirPath, err == nil
}
