
variable "location" {  default = "world" }

module "simple" {
  source = "test-fixtures/simple"
  location = "${var.location}"
}








