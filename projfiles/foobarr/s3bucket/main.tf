

resource "aws_s3_bucket" "default" {
    bucket = "myawesomebucket/infra.backp"
    acl = "private"
    
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
