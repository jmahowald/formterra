package tfproject

import (
	"html/template"
	"io"
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

func getString(key string) string {
	return TfConfig.GetString(key)
}

func getCalledArgs() []string {
	return os.Args
}

var funcMap = template.FuncMap{
	"IsSet":     isSet,
	"GetString": getString,
	"Version":   formTerraVersion,
	"CLIArgs":   getCalledArgs,
}

// TODO this seems like it's a base interace
type templateIterator interface {
	loadTemplate(name string) ([]byte, error)
	getTemplates() ([]string, error)
}

type writerFactory interface {
	getWriter(string) (io.Writer, error)
}

type fileWriter struct {
	targetDir string
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

func (f fileWriter) getWriter(name string) (file io.Writer, err error) {
	_, err = os.Stat(f.targetDir)
	if err != nil {
		log.Warnf("No directory at %s:%s", f.targetDir, err)
		os.MkdirAll(f.targetDir, 0755)
	}
	//Now open the file to write it to and execute
	destFileName := filepath.Join(f.targetDir, name)
	file, err = os.Create(destFileName)
	if err != nil {
		log.Warnf("Unable to create file at %s for writing:%v", destFileName, err)
		return
	}
	return
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
	targetWriter := fileWriter{targetdir}
	for _, dir := range dirs {
		tpls := AssetTemplates{dir}
		err := processTemplates(targetWriter, tpls, context)
		if err != nil {
			log.Warnf("Couldn't process %s", tpls, err)
			return err
		}
	}
	return nil
}

// ProcessTemplates will go through all the templates in a given directory and execute them
func processTemplates(writerFactory writerFactory, templater templateIterator, context TemplateContext) error {

	templates, err := templater.getTemplates()
	if err != nil {
		log.Warnf("Could not find templates with  in %s :%v", templater, err)
		return err
	}
	for _, name := range templates {
		// First load the template and parse it
		templateBytes, err := templater.loadTemplate(name)
		if err != nil {
			log.Fatalf("Unable to retrieve template file %s:%v", name, err)
			return err
		}
		tmpl := template.Must(template.New(name).Funcs(funcMap).Parse(string(templateBytes)))
		log.Debugf("Processing template:%s", name)

		file, err := writerFactory.getWriter(name)
		if err != nil {
			log.Warnf("Unable to create file for %s for writing:%s", name, err)
			return err
		}
		log.Debugf("Going to write into %s into %v ", file, context.getData())
		if err := tmpl.Execute(file, context.getData()); err != nil {
			log.Warnf("Unable to process template %s:%s", name, err)
			return err
		}

	}
	return nil
}
