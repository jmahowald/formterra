package tfproject

// RDSRequest simple request to get a
type RDSRequest struct {
	DatabaseName string
}

// Create Creates a terraform layer to create an s3 bucket
func (req RDSRequest) Create() (TerraformLayer, bool) {

	layer := TerraformLayer{Name: "db" + req.DatabaseName}
	path, _ := layer.getDir()
	processAssetTemplates(path, []string{"rds", "common"}, req)

	return layer, true
}

// GetData returns itself for template contexts
func (s RDSRequest) getData() interface{} {
	return s
}
