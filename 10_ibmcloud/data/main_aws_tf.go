package data

// TerraformAWSTmpl is the Terraform code in form of Go template to be applied by
// Terranova after rendere it
const TerraformAWSTmpl = `
{{ if .SSHAccess }}
variable "public_key_file"  { }
locals {
	public_key    = file(pathexpand(var.public_key_file))
}
{{- end }}

output "public_ip" {
  value   = aws_instance.letschat_server.public_ip
}

provider "aws" {
  region   = "us-west-2"
}

resource "aws_instance" "letschat_server" {
  ami                    = "ami-0d1cd67c26f5fca19"
	instance_type          = "t2.micro"
	key_name      				 = "letschat_server_key"
  vpc_security_group_ids = [ aws_security_group.letschat_server_ingress.id ]
  user_data              = templatefile("${path.module}/user_data.sh", {
    letschat_port       = {{ .LetsChatPort }},
    status              = {{ .Status }},
    status_port         = {{ .StatusPort }},
    docker_compose_b64  = base64encode(templatefile("${path.module}/docker-compose.yaml",{ letschat_port = {{ .LetsChatPort }} })),
  })
	// user_data              = data.template_file.user_data.rendered

  tags = {
    Name = "terranova-example-letschat_server"
	}
}

{{ if .SSHAccess }}
resource "aws_key_pair" "keypair" {
	key_name    = "letschat_server_key"
	public_key  = local.public_key
}
{{- end }}

resource "aws_security_group" "letschat_server_ingress" {
  name = "terranova_letschat_server_ingress"

  ingress {
		description = "Open port to the LetsChat application hosted"
    from_port   = "{{ .LetsChatPort }}"
    to_port     = "{{ .LetsChatPort }}"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
	}

  {{ if .Status }}
  ingress {
		description = "Open port to the LetsChat application hosted"
    from_port   = "{{ .StatusPort }}"
    to_port     = "{{ .StatusPort }}"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  {{- end }}

  {{ if .SSHAccess }}
  ingress {
		description = "Open port to SSH access"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  {{- end }}
	
	egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
`
