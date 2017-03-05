
variable "location" {  default = "world" }
variable "values" {  
  type="list" 
  default = [ 
    "val1", 
    "val2",
  ]
}
  
variable "rds_password" { }
module "simple" {
  source = "test-fixtures/simple"
  location = "${var.location}"
  values = "${var.values}"
  password = "${var.rds_password}"
}






