package data

// TfMainTmplData contain the data to render the template
type TfMainTmplData struct {
	Module string
}

// TfMainTmpl is the template for main.tf
const TfMainTmpl = `
variable "public_key_file" {}
locals {
  public_key = file(pathexpand(var.public_key_file))
}

module "terranova_vm" {
  source 		 = "./modules/{{ .Module }}"
  public_key = local.public_key
}

output "public_ip" {
  value = module.terranova_vm.public_ip
}
`
