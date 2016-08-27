package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"text/template"
)

func parseTemplate(name, path string) *template.Template {
	// var templateBytes byte[]
	var templateBytes []byte
	var err error
	if EmbeddedTemplatesDir == "" {
		templateBytes, err = Asset(fmt.Sprintf("assets/%s", path))
		if err != nil {
			log.Fatalln("Unable to retrieve asset file", err)
		}
	} else {
		templateBytes, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", EmbeddedTemplatesDir, path))
		if err != nil {
			log.Fatalf("Unable to load template file %s/%s:%v", EmbeddedTemplatesDir, path, err)
		}
	}
	tmpl, err := template.New(name).Parse(string(templateBytes))
	if err != nil {
		log.Fatalln("Unable to parse template", err)
	}
	return tmpl
}

const dirMode = 0755
