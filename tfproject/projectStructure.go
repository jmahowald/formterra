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
	VarName      string
	VarValuePath []string
	Type         string
}

//TerraformModuleDefinition Base object for working
//with terraform projects
type TerraformModuleDefinition struct {
	Name          string
	RequiredVars  []string `json:"required_vars"`
	OptionalVars  []string `json:"optional_vars"`
	Outputs       []string
	URI           string
	localLocation string
}

// ModuleCall represents a terraform use of a module
// variables to the modules can come from three types of
// sources, from other modules, from remote sources
// or from variables within the call ing project
type ModuleCall struct {
	TerraformModuleDefinition
	ModuleVariables []FromModuleMappings  `json:"module_vars"`
	RemoteVariables []FromRemoteMappings  `json:"remote_source_vars"`
	Variables       BasicVariableMappings `json:"vars"`
}

//FromModuleMappings variables can come from other
//modules
type FromModuleMappings struct {
	ModuleName string                 `json:"module_name"`
	Mappings   []BasicVariableMapping `json:"mappings"`
}

//FromRemoteMappings variables can come from a terraform
// remote data source
type FromRemoteMappings struct {
	RemoteSourceName string                 `json:"source_name"`
	Mappings         []BasicVariableMapping `json:"mappings"`
	Config           map[string]string      `json:"config,omitempty"`
}

//BasicVariableMappings a list of variables to map in.
type BasicVariableMappings []BasicVariableMapping

//BasicVariableMapping takes input from the source and
//maps it into the given Variable Name for the target
type BasicVariableMapping struct {
	VarName       string `json:"var_name"`
	SourceVarName string `json:"source_var_name,omitempty"`
	Type          string `json:"type,omitempty"`
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
	//TODO this hasn't really been used. Do away with it?
	Name string
}
