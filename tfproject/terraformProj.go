package tfproject

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

const dirMode = 0755

// Overwrite flag setting in config on whether we default overwrite files if they
// exist
const Overwrite = "overwrite"

// TerraformGeneratedProject indicates a terraform project that this
// tool has generated (which means we can safely write in our own contents)
type TerraformGeneratedProject interface {
	Write(data interface{}) (bool, error)
}

// TfConfig configuration for this project.  By default uses viper
var TfConfig = viperConfig

// TerraformLayer is usually a terraform "project", but that is built
// on another layer
type TerraformLayer struct {
	Name         string
	RequiredVars []string
	NextLayer    *TerraformLayer
	// SourceURI    string
}

// PredefinedTerraformProjects terraform projects that are
// embedded within this tool
type PredefinedTerraformProjects struct {
	TerraformLayer
	Templates []string
}

func (t *PredefinedTerraformProjects) Write(data interface{}) {
	t.write()
	for index, element := range t.Templates {
		log.Debug("Processing template:", element)
		tpl := parseTemplate(t.Name+"."+string(index), element)
		f := t.openForWrite(element)
		if err := tpl.Execute(f, data); err != nil {
			log.Fatalln("Unable to generate template:", element, err)
		}
	}

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
}

func (t *TerraformLayer) write() bool {
	dir, isNew := t.dir()
	if isNew || TfConfig.IsSet(Overwrite) {
		log.Debug("Writing out project to ", dir)
		t.makeMake()
		return true
	}
	log.Info("Directory exists, and overwrite is not set")
	return false
}

func (t *TerraformLayer) dir() (string, bool) {
	dirPath := filepath.Join(TfConfig.GetString("terraform-dir"),
		TfConfig.GetString("env"), t.Name)

	_, err := os.Stat(dirPath)
	if err != nil {
		log.Debug("No existing project, generating dir for", dirPath)
		err := os.MkdirAll(dirPath, dirMode)
		if err != nil {
			log.Fatalf("Unable to create directory at %s:%s", dirPath, err)
		}
		return dirPath, true
	}
	return dirPath, false
}

func (t *TerraformLayer) makeMake() bool {
	dir, _ := t.dir()
	_, err := os.Stat(filepath.Join(dir, "Makefile"))
	if err != nil || TfConfig.IsSet(Overwrite) {
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
