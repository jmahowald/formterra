package tfproject

import (
	"os"
	"path"

	log "github.com/Sirupsen/logrus"

	getter "github.com/hashicorp/go-getter"
)

type ExternalModule struct {
	Name string
	Uri  string
}

var externaldirname = "external"

func externaldir() string {
	return path.Join(getString(TerraformDir), externaldirname)
}
func (m ExternalModule) fetch() (TerraformProjectDefinition, error) {
	projectDef := TerraformProjectDefinition{}
	log.Debug("Attempting to retrieve:", m)

	wd, err := os.Getwd()
	if err != nil {
		log.Warn("Couldn't get current working directory")
		return projectDef, err
	}
	srcURI, err := getter.Detect(m.Uri, wd, getter.Detectors)
	if err != nil {
		log.Warn("Could not detect location of", m.Uri, err)
	}
	log.Debug("source uri is ", srcURI)

	if projectDef.Name != "" {
		projectDef.Name = path.Base(srcURI)
	}

	destPath := path.Join(externaldir(), projectDef.Name)
	err = getter.Get(destPath, srcURI)
	if err != nil {
		log.Warn("Unable to retrieve:", srcURI)
		return projectDef, err
	}
	projectDef.location = destPath
	log.Debug("Retrieved:", projectDef)
	projectDef.loadVars()
	return projectDef, nil
}

func (t TerraformProjectDefinition) loadVars() {

	t.RequiredVars = []string{"hello"}
}

func init() {
	os.MkdirAll(externaldir(), 0755)
}
