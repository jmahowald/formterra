variable "subnet_ids" {
  type = "list"
}

variable "vpc_id"{}
variable "inbound_security_group_ids" {}

variable "username" {
  default = "root"
}
variable "multi_az" {
  default = true
}
variable "password" {}
variable "schema_name"{}
variable "engine" {
  default = "mysql"
}
variable "engine_version" {
  default = "5.6.23"
}
variable "instance_class" {
  default = "db.t2.micro"
}


# Database instance
resource "aws_db_instance" "default" {
    allocated_storage = 10
    engine = "${var.engine}"
    engine_version = "${var.engine_version}"
    identifier = "{{.DBName}}"
    instance_class = "${var.instance_class}"
    //TODO figure out what the heck this means
    final_snapshot_identifier = "{{.DBName}}-final"
    publicly_accessible = false
    db_subnet_group_name = "${var.db_subnet_group_name}"
    vpc_security_group_ids = [
        "${var.db_security_group}"
    ]
    multi_az = "${var.multi_az}"

    # Database details
    name = "{{.DBName}}"
    username = "${var.username}"
    password = "${var.password}"

    lifecycle {
        create_before_destroy = true
    }
}

output "database_address" {
  value = "${aws_db_instance.default.address}"
}
