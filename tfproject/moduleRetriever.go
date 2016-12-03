package tfproject

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	// set "github.com/deckarep/golang-set"
	getter "github.com/hashicorp/go-getter"
)

// This may well be overkill
//ExternalModule  a terraform module that isn't baked into this project
type ExternalModule struct {
	Name string
	URI  string
}

var externaldirname = "external"

func externaldir() string {
	return path.Join(getString(TerraformDir), externaldirname)
}

// grabs the remote definition (which might be local)
func (m ExternalModule) Fetch() (TerraformModuleDefinition, error) {
	moduleDef := TerraformModuleDefinition{URI: m.URI}

	wd, err := os.Getwd()
	if err != nil {
		log.Warn("Couldn't get current working directory")
		return moduleDef, err
	}

	baseUri, subdir := getter.SourceDirSubdir(m.URI)

	srcURI, err := getter.Detect(baseUri, wd, getter.Detectors)
	log.Debugf("base uri is %s subdir is %s", srcURI, subdir)

	if err != nil {
		log.Warnf("Could not detect location of %s:%s", m.URI, err)
		return moduleDef, err
	}

	if moduleDef.Name == "" {
		moduleDef.Name = path.Base(m.URI)
	}

	destPath := path.Join(externaldir(), path.Base(baseUri))
	err = getter.Get(destPath, srcURI)
	if err != nil {
		log.Warn("Unable to retrieve:", srcURI)
		return moduleDef, err
	}
	moduleDef.URI = m.URI

	// localLocation, err := filepath.Abs(filepath.Join(destPath, subdir))
	// if err !
	moduleDef.LocalLocation = filepath.Join(destPath, subdir)

	moduleDef.loadVars()
	return moduleDef, nil
}

type terraformvars []map[string][]map[string]interface{}
type projectVars struct {
	required []string
	optional []string
}

// So while viper can parse out "variables" from terraform, it is
// a little bit of a mess. All the "variable" entries are a map of variable
// names with the values being a list of maps, for each "key".  in the case of
// properly done terraform that is either a description or a key for the "default"
// value, but it also appears to allow other valid keys as well.
// We are interested in detecting the presence of a "default" value for a
// variable, not caring whether that value is a string, list, or a map
func findRequiredAndOptionalVars(vars terraformvars) projectVars {
	//The number of required variables or optional can't be more than
	//the total number. Just go ahead and allocate that much into a slice
	required := make([]string, 0, len(vars))
	optional := make([]string, 0, len(vars))

	for _, varentry := range vars {
		for varname, varkeys := range varentry {
			log.Debug("Variable entry,", varname, ",", varkeys)
			if len(varkeys) == 0 && len(varname) != 0 {
				log.Debug("Found variable without any keys, adding to required:", varname)
				required = append(required, varname)
				//If there are no "keys" under a variable,
				//there is no default
			} else {
				for _, key := range varkeys {
					if _, hasDefault := key["default"]; hasDefault {
						log.Debug("Found variable with default key:", varname)
						optional = append(optional, varname)
					} else {
						log.Debug("Found variable without key:", varname)
						required = append(required, varname)
					}
				}
			}
		}
	}

	varDef := projectVars{required, optional}
	return varDef
}

//
// type VarStruct struct {
// 		defs string,
// }

func (t *TerraformModuleDefinition) loadVars() error {

	files, _ := filepath.Glob(fmt.Sprintf("%s/*.tf", t.LocalLocation))

	for _, tfFile := range files {
		varConfig := viper.New()

		log.Debug("parsing %s", tfFile)
		varConfig.SetConfigFile(tfFile)
		varConfig.SetConfigType("hcl")
		varConfig.ReadInConfig()

		vars := varConfig.Get("variable")

		var variableContents terraformvars
		mapstructure.Decode(vars, &variableContents)
		varLists := findRequiredAndOptionalVars(variableContents)

		log.Debug("required vars are:", varLists.required)
		log.Debug("optional vars are:", varLists.optional)

		t.RequiredVars = append(t.RequiredVars, varLists.required...)
		t.OptionalVars = append(t.OptionalVars, varLists.optional...)

		var outputs terraformvars

		outMap := varConfig.Get("output")
		mapstructure.Decode(outMap, &outputs)
		t.Outputs = append(t.Outputs, findOutputs(outputs)...)

	}

	return nil
}

func findOutputs(vars terraformvars) []string {
	outputs := make([]string, 0, len(vars))
	for _, varentry := range vars {
		for varname, _ := range varentry {
			outputs = append(outputs, varname)
		}
	}
	return outputs
}

func init() {
	os.MkdirAll(externaldir(), 0755)
}
