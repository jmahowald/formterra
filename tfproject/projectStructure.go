package tfproject

import "os/exec"

// MappingType variables can come from tfvars, from another module or
// from a remote location
type MappingType int

// types of variable mappings
const (
	TFVARS MappingType = iota
	MODULE
	REMOTE
)

// VarMappings convenience method
type VarMappings []VarMapping

// VarMapping maps a required variable to a source for terraform
type VarMapping struct {
	VarName        string
	VarMappingType MappingType
	VarValuePath   []string
}

//TerraformProjectDefinition Base object for working
//with terraform projects
type TerraformProjectDefinition struct {
	Name          string
	RequiredVars  []string
	OptionalVars  []string
	URI           string
	localLocation string
}

type FromVariableMapping struct {
	VarName string
}

type FromModuleMappings struct {
	ModuleName string                 `json:"module_name"`
	Mappings   []BasicVariableMapping `json:"mappings"`
}

type FromRemoteMappings struct {
	RemoteSourceName string                 `json:"source_name"`
	Mappings         []BasicVariableMapping `json:"mappings"`
}

type BasicVariableMapping struct {
	VarName       string `json:"var_name"`
	SourceVarName string `json:"source_var_name"`
}

type ModuleCall struct {
	URI             string `json:"uri"`
	Name            string
	ModuleVariables []FromModuleMappings   `json:"module_vars"`
	RemoteVariables []FromRemoteMappings   `json:"remote_source_vars"`
	Variables       []BasicVariableMapping `json:"vars"`
}

type ModuleTwo struct {
	URI              string
	Name             string
	ModulesVariables []struct {
		ModuleName string
		Variables  []struct {
			OutputName   string
			VariableName string
		}
	}
}

// TerraformProjectSkeleton a terraform project
// A terraform project will typically
// call off to many modules
type TerraformProjectSkeleton struct {
	TerraformLayer
	Modules []ModuleCall
}

// TerraformProject something that we can see the plan for, apply, and destroy
type TerraformProject interface {
	Plan() exec.Cmd
	Apply() exec.Cmd
	Destroy() exec.Cmd
}

// TerraformLayer is usually a terraform "project", but that is built
// on another layer
// The layers are stored in a consistent directory
type TerraformLayer struct {
	Name string
	// SourceURI    string
}

//BuiltInTerraformProjectRequest we already have the defintiion of the templates
//because they are embedded in the tool
type TemplateRequest struct {
	name      string
	templates []string
	data      interface{}
}
