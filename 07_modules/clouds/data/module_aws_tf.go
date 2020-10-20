package data

// TfModAwsMainCode is the plain text code for module/aws/main.tf
const TfModAwsMainCode = `
variable "public_key" {}

provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "terranova_vm" {
  instance_type = "t2.micro"
  ami           = "ami-6e1a0117"
  key_name      = "terranovaKeyPair"
}

resource "aws_key_pair" "terranova_key_pair" {
  key_name   = "terranovaKeyPair"
  public_key = var.public_key
}

output "public_key" {
  value = aws_instance.terranova_vm.public_ip
}
`
