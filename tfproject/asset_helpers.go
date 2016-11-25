package tfproject

import (
	"html/template"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type TemplateContext interface {
	getData() interface{}
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
	"Version":   formTerraVersion,
}

// TODO this seems like it's a base interace
type templateIterator interface {
	loadTemplate(name string) ([]byte, error)
	getTemplates() ([]string, error)
}

type AssetTemplates struct {
	assetFolder string
}

type DirTemplates struct {
	dirName string
}

func (at AssetTemplates) getTemplates() ([]string, error) {
	templates, err := AssetDir(filepath.Join("assets", at.assetFolder))
	if err != nil {
		log.Fatalf("Couldn't load any assets from %s:%v", at.assetFolder, err)
		return templates, err
	}
	return templates, nil
}

func (at AssetTemplates) loadTemplate(name string) ([]byte, error) {
	return Asset(filepath.Join("assets", at.assetFolder, name))
}

// GenerateSkeleton creates a terraform project
func (t TerraformProjectSkeleton) GenerateSkeleton() error {
	//For now ignore if the directory already exists
	dir, _ := t.getDir()
	err := processAssetTemplates(dir, []string{"project", "common"}, t)
	return err
}

func (t TerraformProjectSkeleton) getData() interface{} {
	return t
}

func processAssetTemplates(targetdir string, dirs []string, context TemplateContext) error {
	for _, dir := range dirs {
		tpls := AssetTemplates{dir}
		err := ProcessTemplates(targetdir, tpls, context)
		if err != nil {
			log.Warnf("Couldn't process %s", tpls, err)
			return err
		}
	}
	return nil
}

// ProcessTemplates will go through all the templates in a given directory and execute them
func ProcessTemplates(targetdir string, templater templateIterator, context TemplateContext) error {
	log.Debugf("Processings templates with %s to be placed in %s", templater, targetdir)
	_, err := os.Stat(targetdir)
	if err != nil {
		log.Warnf("No directory at %s:%s", targetdir, err)
		os.MkdirAll(targetdir, 0755)
	}

	templates, err := templater.getTemplates()
	if err != nil {
		log.Warnf("Could not find templates with  in %s :%v", templater, err)
		return err
	}
	for _, name := range templates {
		// First load the template and parse it
		log.Debugf("Processing template:%s", name)
		templateBytes, err := templater.loadTemplate(name)
		if err != nil {
			log.Fatalf("Unable to retrieve template file %s:%v", name, err)
			return err
		}
		tmpl := template.Must(template.New(name).Funcs(funcMap).Parse(string(templateBytes)))

		//Now open the file to write it to and execute
		destFileName := filepath.Join(targetdir, name)
		file, err := os.Create(destFileName)
		if err != nil {
			log.Warnf("Unable to create file at %s for writing:%v", destFileName, err)
			return err
		}
		if err := tmpl.Execute(file, context.getData()); err != nil {
			log.Warnf("Unable to process template %s:%s", name, err)
			return err
		}
	}
	return nil
}
