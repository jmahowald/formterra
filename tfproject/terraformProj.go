package tfproject

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"

	"github.com/jmahowald/formterra/core"
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
	log.Debug("Attempting to create:", t.Name)
	fullPath, exists := layerExists(t.Name)
	if !exists {
		err := os.MkdirAll(fullPath, dirMode)
		if err != nil {
			log.Fatalf("Unable to create directory at %s:%s", fullPath, err)
		}
	}
	return fullPath, exists
}

//Create builds layer by executing all of the templates
func (req TemplateRequest) Create() (TerraformLayer, bool) {
	layer := TerraformLayer{Name: req.name}
	fullPath, exists := layer.getDir()
	if !exists || isSet(Overwrite) {
		log.Debug("Writing out project to ", fullPath)
		layer.makeMake()
		for index, element := range req.templates {
			log.Debug("Processing template:", element)
			tpl := parseTemplate(layer.Name+"."+string(index), element)
			f := layer.openForWrite(element)
			if err := tpl.Execute(f, req.data); err != nil {
				log.Fatalln("Unable to generate template:", element, err)
			}
		}
	} else {
		log.Info("Directory exists, and overwrite is not set")
	}
	return layer, exists
}

// ParseTemplate
func parseTemplate(name, path string) *template.Template {
	// var templateBytes byte[]
	var templateBytes []byte
	templateBytes = loadAsset(path)

	tmpl := template.Must(template.New(name).Funcs(funcMap).Parse(string(templateBytes)))
	return tmpl
}

//IsSet used by templates to see if a value is set
func isSet(key string) bool {
	return TfConfig.IsSet(key)
}

//GetSt
func getString(key string) string {
	return TfConfig.GetString(key)
}

var funcMap = template.FuncMap{
	"IsSet":     isSet,
	"GetString": getString,
	"version":   formTerraVersion,
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

func (t *TerraformLayer) makeMake() bool {
	dir, _ := t.dir()
	_, err := os.Stat(filepath.Join(dir, "Makefile"))
	if err != nil || isSet(Overwrite) {
		makeContents, err := Asset(filepath.Join("assets", "Makefile"))
		if err != nil {
			log.Fatalln("Unable to retrieve base Makefile file", err)
		}
		err = ioutil.WriteFile(filepath.Join(dir, "Makefile"), makeContents, 0644)
		if err != nil {
			log.Fatalln("Unable to create Makefile", err)
		}
		return true
	}
	return false
}

func (t *TerraformLayer) openForWrite(fileName string) *os.File {
	dir, _ := t.dir()
	file, err := os.Create(filepath.Join(dir, fileName))
	//TODO test overwrite?
	if err != nil {
		log.Fatalf("Unable to create file at %s:%v", fileName, err)
	}
	return file
}
