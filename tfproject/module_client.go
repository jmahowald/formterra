package tfproject

import (
	"strings"

	log "github.com/Sirupsen/logrus"
)

//TerraformProjectDefinition Base object for working with terraform projects
// type TerraformProjectDefinition struct {
// 	Name         string
// 	RequiredVars []string
// 	OptionalVars []string
// 	location     string
// }
//
// func (t *TerraformProjectDefinition) CreateClient() error {
// 	return nil
// }

type varMap map[string]string

type moduleRequest struct {
	Name      string
	SourceURI string
	Mappings  VarMappings
}

// VarMapFactory simple interface to allow us to
// easily create a call off to a module
type VarMapFactory interface {
	CreateMapping() VarMapping
}

// VarPath Get the formatted path to the value for terraform
// interpolation
func (v VarMapping) VarPath() string {
	return strings.Join(v.VarValuePath, ".")
}

func generateModule(moduleURI, name string) {

	mod := ExternalModule{name, moduleURI}
	def, err := mod.Fetch()
	if err != nil {
		log.Fatalf("could not fetch module from %s", moduleURI)
	}

	mappings := simpleMappings(def.RequiredVars)
	modDef := moduleRequest{def.Name, moduleURI, mappings}
	req := TemplateRequest{name,
		[]string{"module_client.tf"},
		modDef}

	req.Create()
}

func simpleMappings(requiredVars []string) VarMappings {
	mappings := make([]VarMapping, len(requiredVars), len(requiredVars))
	for i := range requiredVars {
		mappings[i] = newVariableMapping(requiredVars[i], requiredVars[i])
	}
	return mappings
}

// NewModuleMapping a mapping of a variable to
// to another modules source
func newModuleMapping(varName, moduleName, mappedValue string) VarMapping {
	varPath := []string{"module", moduleName, mappedValue}
	x := VarMapping{
		VarName:        varName,
		VarMappingType: MODULE,
		VarValuePath:   varPath,
	}
	return x
}

// NewVariableMapping A mapping that should come from a tfvars file or
// environment variable
func newVariableMapping(varName, mappedValue string) VarMapping {
	varPath := []string{"var", mappedValue}
	x := VarMapping{
		VarName:        varName,
		VarMappingType: TFVARS,
		VarValuePath:   varPath,
	}
	return x
}

// NewRemoteStateMapping creates a mapping for variables
func newRemoteStateMapping(varName, stateName, mappedValue string) VarMapping {
	varPath := []string{"data", "terraform_remote_state", stateName, mappedValue}
	x := VarMapping{
		VarName:        varName,
		VarMappingType: REMOTE,
		VarValuePath:   varPath,
	}
	return x
}
