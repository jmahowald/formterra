package tfproject

// S3BucketRequest having a base fqdn and a name on top of that.
type S3BucketRequest struct {
	Fqdn        string
	BucketName  string
	UnVersioned bool
}

// TerraformS3  creates a terraform project from a bucket request
func TerraformS3(bucketRequest S3BucketRequest) TerraformGeneratedProject {
	layer := TerraformLayer{Name: "bucket"}
	s3Proj := TemplatedTerraformProjects{
		TerraformLayer: layer,
		Templates:      []string{"s3.tf"},
	}
	return s3Proj
}
