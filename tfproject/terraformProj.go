package tfproject

import (
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

const dirMode = 0755

// TerraformGeneratedProject indicates a terraform project that this
// tool has generated (which means we can safely write in our own contents)
type TerraformGeneratedProject interface {
	Write(data interface{})
}

// TerraformProject something that we can see the plan for, apply, and destroy
type TerraformProject interface {
	Plan() exec.Cmd
	Apply() exec.Cmd
	Destroy() exec.Cmd
}

type makeable interface {
}

type terraformProjectRequest interface {
	Create() (TerraformLayer, bool)
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

// type projectRequest struct {
// 	data interface{}
// 	TerraformGeneratedProject
// }

type BuiltInTerraformProjectRequest struct {
	name      string
	templates []string
	data      interface{}
}

// TemplatedTerraformProjects terraform projects that are
// embedded within this tool
type TemplatedTerraformProjects struct {
	TerraformLayer
	Templates []string
}

func (req BuiltInTerraformProjectRequest) Create() (TerraformLayer, bool) {
	log.Debug("Attempting to create:", req)
	fullPath, exists := layerExists(req.name)
	if !exists {
		err := os.MkdirAll(fullPath, dirMode)
		if err != nil {
			log.Fatalf("Unable to create directory at %s:%s", fullPath, err)
		}
	}

	layer := TerraformLayer{Name: req.name}
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
}

//
// func command(t *TerraformLayer, command string, args ...string) (exec.Cmd, error) {
// 	bin, err := exec.LookPath(command)
// 	if err != nil {
// 		log.Warn("Could not find in path:", command, err)
// 		return exec.Cmd{}, err
// 	}
// 	dir, _ := t.dir()
// 	cmd := exec.Cmd{Path: bin, Dir: dir, Args: args}
// 	return cmd, nil
// }

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
