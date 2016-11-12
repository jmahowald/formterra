package tfproject

import "strings"

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

// VarPath Get the formatted path to the value for terraform
// interpolation
func (v VarMapping) VarPath() string {
	return strings.Join(v.VarValuePath, ".")
}

func (v BasicVariableMapping) interpolationPath(prefix []string) VarMapping {
	var sourceName string
	if v.SourceVarName != "" {
		sourceName = v.SourceVarName
	} else {
		sourceName = v.VarName
	}
	return VarMapping{v.VarName, append(prefix, sourceName)}
}

func (m ModuleCall) GetVariables() VarMappings {
	// Default of 20 mappings. No good reason for this
	// heuristic. because we use append we don't really
	// worry about this being to small
	mappings := make([]VarMapping, 0, 20)
	for _, modVars := range m.ModuleVariables {
		mappings = append(mappings, modVars.GetTerraformMappings()...)
	}
	for _, remoteVars := range m.RemoteVariables {
		mappings = append(mappings, remoteVars.GetTerraformMappings()...)
	}
	mappings = append(mappings, m.Variables.GetTerraformMappings()...)
	return mappings
}

func (t TerraformProjectSkeleton) GetAllVars() VarMappings {
	// Default of 20 mappings. No good reason for this
	// heuristic. because we use append we don't really
	// worry about this being to small
	mappings := make([]VarMapping, 0, 20)
	for _, module := range t.Modules {
		mappings = append(mappings, module.Variables.GetTerraformMappings()...)
	}
	return mappings
}

//
// func generateModule(moduleURI, name string) {
//
// 	mod := ExternalModule{name, moduleURI}
// 	def, err := mod.Fetch()
// 	if err != nil {
// 		log.Fatalf("could not fetch module from %s", moduleURI)
// 	}
//
// 	mappings := simpleMappings(def.RequiredVars)
// 	modDef := moduleRequest{def.Name, moduleURI, mappings}
// 	req := TemplateRequest{name,
// 		[]string{"module_client.tf"},
// 		modDef}
//
// 	req.Create()
// }

//variableSourceMapper able to get terraform interpolations
type variableSourceMapper interface {
	GetTerraformMappings() VarMappings
}

func (v BasicVariableMappings) GetTerraformMappings() VarMappings {
	mappings := make([]VarMapping, len(v), len(v))
	prefix := []string{"var"}
	for i := range v {
		mappings[i] = v[i].interpolationPath(prefix)
	}
	return mappings
}

func (mod FromModuleMappings) GetTerraformMappings() VarMappings {
	vars := mod.Mappings
	prefix := []string{"module", mod.ModuleName}
	mappings := make([]VarMapping, len(vars), len(vars))
	for i := range vars {
		mappings[i] = vars[i].interpolationPath(prefix)
	}
	return mappings
}

func (remote FromRemoteMappings) GetTerraformMappings() VarMappings {
	prefix := []string{"data", "terraform_remote_state", remote.RemoteSourceName}
	vars := remote.Mappings
	mappings := make([]VarMapping, len(vars), len(vars))
	for i := range vars {
		mappings[i] = vars[i].interpolationPath(prefix)
	}
	return mappings
}
