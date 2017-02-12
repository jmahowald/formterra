variable "greeting" {
  default = "hello"
}

variable "location" {}

variable "values" {
  type="list"
  default=[]
}
output "out1" {
  value = "${var.greeting}"
}
output "out2" {
  value = "${var.greeting}"
}