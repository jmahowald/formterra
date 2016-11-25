variable "greeting" {
  default = "hello"
}

variable "location"{}

output "out1" {
  value = "${var.greeting}"
}

output "out2" {
  value = "${var.greeting}"
}