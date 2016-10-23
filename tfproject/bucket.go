package tfproject

type S3BucketID struct {
	Fqdn       string
	BucketName string
}

// S3BucketRequest having a base fqdn and a name on top of that.
type S3BucketRequest struct {
	S3BucketID
	UnVersioned bool
}

// Create Creates a terraform layer to create an s3 bucket
func (s S3BucketRequest) Create() (TerraformLayer, bool) {
	request := BuiltInTerraformProjectRequest{
		name:      "bucket_" + s.BucketName,
		templates: []string{"s3.tf"},
		data:      s,
	}
	return request.Create()
}
