

resource "aws_s3_bucket" "default" {
    bucket = "{{.BucketName}}/{{.Fqdn}}"
    acl = "private"
    {{  if not .UnVersioned }}
      versioning {
        enabled = true
      }
      /**
       * Dont' keep around all old versions forever
       * TODO make this configurable
       */
      lifecycle_rule {
          prefix = "config/"
          enabled = true
          noncurrent_version_transition {
              days = 30
              storage_class = "STANDARD_IA"
          }
          noncurrent_version_transition {
              days = 60
              storage_class = "GLACIER"
          }
          noncurrent_version_expiration {
              days = 90
          }
      }
      
    {{ end }}
}


output "bucket_region" {
  value = "${aws_s3_bucket.default.region}"
}
output "arn" {
  value = "${aws_s3_bucket.default.arn}"
}
output "name" {
  value = "${aws_s3_bucket.default.id}"
}
