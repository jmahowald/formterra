data "template_file" "s3-ro" {
  template = "${file("policy.s3-ro.tpl")}"
  vars {
    bucket_name= "${var.bucket_name}"
  }
}


data "template_file" "s3-rw" {
  template = "${file("policy.s3-rw.tpl")}"
  vars {
    bucket_name= "${var.bucket_name}"
  }
}

data "template_file" "s3-principal" {
  template = "${file("policy.s3-principal.tpl")}"
  vars {
    bucket_name= "${var.bucket_name}"
  }
}


