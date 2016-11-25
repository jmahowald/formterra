package tfproject

// TODO make sure name doesn't have periods in it
type S3BucketID struct {
	Fqdn       string
	BucketName string
}

// S3BucketRequest having a base fqdn and a name on top of that.
type S3BucketRequest struct {
	S3BucketID
	UnVersioned bool
	CreateUser  bool
}

// Create Creates a terraform layer to create an s3 bucket
func (s S3BucketRequest) Create() (TerraformLayer, bool) {

	layer := TerraformLayer{Name: "bucket_" + s.BucketName}
	path, _ := layer.getDir()
	processAssetTemplates(path, []string{"s3", "common"}, s)

	return layer, true
}

// GetData returns itself for template contexts
func (s S3BucketRequest) getData() interface{} {
	return s
}
