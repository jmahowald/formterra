

{{  if .CreateUser }}


resource "aws_iam_user" "{{.BucketName}}" {
    name = "{{.BucketName}}"
    path = "/system/"
}

output "access_key" {
  value = "${aws_iam_access_key.{{.BucketName}}.id}"
}
output "secret" {
  value = "${aws_iam_access_key.{{.BucketName}}.secret}"
}


resource "aws_iam_access_key" "{{.BucketName}}" {
    user = "${aws_iam_user.{{.BucketName}}.name}"
}

resource "aws_iam_user_policy_attachment" "rwpolicy" {
     user = "${aws_iam_user.{{.BucketName}}.name}"
     policy_arn = "${aws_iam_policy.s3-rw.arn}"
}

{{ else }}

// No user requested 

{{  end }}

