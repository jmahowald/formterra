package tfproject

import (
	"log"
	"os/exec"
)

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
}

//BasicVariableMappings a list of variables to map in.
type BasicVariableMappings []BasicVariableMapping

//BasicVariableMapping takes input from the source and
//maps it into the given Variable Name for the target
type BasicVariableMapping struct {
	VarName       string `json:"var_name"`
	SourceVarName string `json:"source_var_name"`
}

// ModuleCall represents a terraform use of a module
// variables to the modules can come from three types of
// sources, from other modules, from remote sources
// or from variables within the call ing project
type ModuleCall struct {
	URI             string `json:"uri"`
	Name            string
	ModuleVariables []FromModuleMappings  `json:"module_vars"`
	RemoteVariables []FromRemoteMappings  `json:"remote_source_vars"`
	Variables       BasicVariableMappings `json:"vars"`
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

//TemplateRequest we already have the defintiion of the templates
//because they are embedded in the tool
type TemplateRequest struct {
	name      string
	templates []string
	data      interface{}
}

func (t TerraformProjectSkeleton) GenerateSkeleton() error {

	tpl := parseTemplate("project", "project.tf")
	moduleTemplateBytes := loadAsset("module_client.tf")
	tpl, err := tpl.Parse(string(moduleTemplateBytes))
	t.getDir()
	f := t.openForWrite("main.tf")
	if err = tpl.Execute(f, t); err != nil {
		log.Fatalln("Unable to execute template", t, err)
		return err
	}
	return nil
}
