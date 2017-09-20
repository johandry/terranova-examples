resource "aws_instance" "example" {
  ami           = "ami-6e1a0117"
  instance_type = "t2.micro"
}
