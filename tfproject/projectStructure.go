package tfproject

type ModuleCall struct {
}

// MappingType variables can come from tfvars, from another module or
// from a remote location
type MappingType int

// types of variable mappings
const (
	TFVARS MappingType = iota
	MODULE
	REMOTE
)

// VarMapping maps a required variable to a source for terraform
type VarMapping struct {
	VarName        string
	VarMappingType MappingType
	VarValuePath   []string
}

// VarMappings convenience method
type VarMappings []VarMapping

//TerraformProjectDefinition Base object for working
//with terraform projects
type TerraformProjectDefinition struct {
	Name         string
	RequiredVars []string
	OptionalVars []string
	location     string
}

type ModuleInfoRetriever interface {
	GetModuleData(name, sourceURI string) (TerraformProjectDefinition, error)
}
