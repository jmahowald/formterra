package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func config(key string) string {
	return viper.GetString(key)
}

// TerraformLayer is usually a terraform "project", but that is built
// on another layer
type TerraformLayer struct {
	RequiredVars []string
	SourceURI    string
	Name         string
}

func (t *TerraformLayer) dir() string {
	dirPath := filepath.Join(config("terraform-dir"), config("env"), t.Name)
	err := os.MkdirAll(dirPath, dirMode)
	if err != nil {
		log.Fatalf("Unable to create directory at %s:%s", dirPath, err)
	}
	return dirPath
}

type makeargs struct {
	TerraformLayer
}

func (t *TerraformLayer) makeMake() bool {
	dir := t.dir()
	_, err := os.Stat(filepath.Join(dir, "Makefile"))
	if err != nil || OverWriteFiles {
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
	dir := t.dir()
	file, err := os.Create(filepath.Join(dir, fileName))
	//TODO test overwrite?
	if err != nil {
		log.Fatalf("Unable to create file at %s:%v", fileName, err)
	}
	return file
}
