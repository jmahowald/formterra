package tfproject

import "strings"

//TerraformModuleDefinition Base object for working with terraform projects
// type TerraformModuleDefinition struct {
// 	Name         string
// 	RequiredVars []string
// 	OptionalVars []string
// 	location     string
// }
//
// func (t *TerraformModuleDefinition) CreateClient() error {
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
	return VarMapping{v.VarName, append(prefix, sourceName), v.Type, v.DefaultValue}
}

func (m ModuleCall) GetVariables() VarMappings {
	// Default of 20 mappings. No good reason for this
	// heuristic. because we use append we don't really
	// worry about this being to small
	mappings := make([]VarMapping, 0, 20)
	for _, modVars := range m.ModuleVariables {
		mappings = append(mappings, modVars.getTerraformMappings()...)
	}
	for _, remoteVars := range m.RemoteVariables {
		mappings = append(mappings, remoteVars.getTerraformMappings()...)
	}
	mappings = append(mappings, m.Variables.getTerraformMappings()...)
	return mappings
}

func (t TerraformProjectSkeleton) GetAllVars() VarMappings {
	// Default of 20 mappings. No good reason for this
	// heuristic. because we use append we don't really
	// worry about this being to small
	encountered := map[string]bool{}

	mappings := make([]VarMapping, 0, 20)
	for _, module := range t.Modules {
		for _, moduleVar := range module.Variables.getTerraformMappings() {
			if encountered[moduleVar.VarName] == true { //duplicate
			} else {
				encountered[moduleVar.VarName] = true
				mappings = append(mappings, moduleVar)
			}
		}
	}
	return mappings
}

//variableSourceMapper able to get terraform interpolations
type variableSourceMapper interface {
	getTerraformMappings() VarMappings
}

func (v BasicVariableMappings) getTerraformMappings() VarMappings {
	mappings := make([]VarMapping, len(v), len(v))
	prefix := []string{"var"}
	for i := range v {
		mappings[i] = v[i].interpolationPath(prefix)
	}
	return mappings
}

func (mod FromModuleMappings) getTerraformMappings() VarMappings {
	vars := mod.Mappings
	prefix := []string{"module", mod.ModuleName}
	mappings := make([]VarMapping, len(vars), len(vars))
	for i := range vars {
		mappings[i] = vars[i].interpolationPath(prefix)
	}
	return mappings
}

func (remote FromRemoteMappings) getTerraformMappings() VarMappings {
	prefix := []string{"data", "terraform_remote_state", remote.RemoteSourceName}
	vars := remote.Mappings
	mappings := make([]VarMapping, len(vars), len(vars))
	for i := range vars {
		mappings[i] = vars[i].interpolationPath(prefix)
	}
	return mappings
}

// CreateSkeleton Transforms a set of modules into a skeleton
func CreateSkeleton(mods []TerraformModuleDefinition, name string) TerraformProjectSkeleton {
	skel := TerraformProjectSkeleton{}
	skel.TerraformLayer = TerraformLayer{Name: name}
	for _, mod := range mods {
		modVars := make([]BasicVariableMapping, len(mod.RequiredVars))
		for i, variable := range mod.RequiredVars {
			modVars[i] = BasicVariableMapping{VarName: variable}
		}
		modCall := ModuleCall{TerraformModuleDefinition: mod, Variables: modVars}
		skel.Modules = append(skel.Modules, modCall)
	}
	return skel
}
