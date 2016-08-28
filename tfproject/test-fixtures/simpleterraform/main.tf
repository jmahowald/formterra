
resource null_resource "simple" {
  triggers {
    greeting = "${var.greeting}"
    location = "${var.location}"
  }
  provisioner "local-exec" {
       command = "echo ${var.greeting} ${var.location} >> msg.txt"
   }
}
