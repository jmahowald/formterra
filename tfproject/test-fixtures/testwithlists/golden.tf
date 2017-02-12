
variable "location" {  default = "world" }
variable "values" {  
  type="list" 
  default = [ 
    "val1", 
    "val2",
  ]
}
  
module "simple" {
  source = "test-fixtures/simple"
  location = "${var.location}"
  values = "${var.values}"
}








