package tfproject

import (
	"os"
	"path"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	getter "github.com/hashicorp/go-getter"
)

//ExternalModule  a terraform module that isn't baked into this project
type ExternalModule struct {
	Name string
	URI  string
}

//type VarMap []map[string]map[string]interface{}

var externaldirname = "external"

func externaldir() string {
	return path.Join(getString(TerraformDir), externaldirname)
}

// grabs the remote definition (which might be local)
func (m ExternalModule) fetch() (TerraformProjectDefinition, error) {
	projectDef := TerraformProjectDefinition{}
	log.Debug("Attempting to retrieve:", m)

	wd, err := os.Getwd()
	if err != nil {
		log.Warn("Couldn't get current working directory")
		return projectDef, err
	}
	srcURI, err := getter.Detect(m.URI, wd, getter.Detectors)
	if err != nil {
		log.Warn("Could not detect location of", m.URI, err)
	}
	log.Debug("source uri is ", srcURI)

	if projectDef.Name == "" {
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
						log.Debug("Found variable with default key", varname)
						optional = append(optional, varname)
					} else {
						log.Debug("Found variable without key", varname)
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

func (t *TerraformProjectDefinition) loadVars() error {

	variablesFile := path.Join(t.location, "variables.tf")
	if _, err := os.Stat(variablesFile); err != nil {
		log.Info("No variables.tf in:", t.location)
		return err
	}

	varConfig := viper.New()
	varConfig.SetConfigFile(variablesFile)
	varConfig.SetConfigType("hcl")
	varConfig.ReadInConfig()
	varConfig.Debug()

	vars := varConfig.Get("variable")
	log.Debug("Vars loaded are:", reflect.TypeOf(vars), vars)

	var variableContents terraformvars
	mapstructure.Decode(vars, &variableContents)
	log.Debug("Variable Contents:", variableContents)
	projectVarDef := findRequiredAndOptionalVars(variableContents)
	t.RequiredVars = projectVarDef.required
	t.OptionalVars = projectVarDef.optional

	return nil
}

func init() {
	os.MkdirAll(externaldir(), 0755)
}