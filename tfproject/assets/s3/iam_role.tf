#resource "aws_iam_instance_profile" "test-profile" {
#  name = "{{.BucketName}}-profile"
#  roles = [
#    "${aws_iam_role.s3-ro.name}",
#    "${aws_iam_role.s3-rw.name}"
#  ]
#}

variable "bucket_name" { default = "{{.BucketName}}-bucket" }

resource "aws_iam_role" "s3-ro" {
  name = "{{.BucketName}}-s3-ro"
  path = "/"
  assume_role_policy = "${template_file.s3-principal.rendered}"
}

resource "aws_iam_role" "s3-rw" {
  name = "{{.BucketName}}-s3-rw"
  path = "/"
  assume_role_policy = "${template_file.s3-principal.rendered}"
}

resource "aws_iam_policy" "s3-ro" {
  name = "{{.BucketName}}-s3-ro"
  policy = "${template_file.s3-ro.rendered}"  
}

resource "aws_iam_policy" "s3-rw" {
  name = "{{.BucketName}}-s3-rw"
  policy = "${template_file.s3-rw.rendered}"
}

resource "aws_iam_role_policy_attachment" "s3-ro" {
  role = "${aws_iam_role.s3-ro.name}"
  policy_arn = "${aws_iam_policy.s3-ro.arn}"
}

resource "aws_iam_role_policy_attachment" "s3-rw" {
  role = "${aws_iam_role.s3-rw.name}"
  policy_arn = "${aws_iam_policy.s3-rw.arn}"
}
