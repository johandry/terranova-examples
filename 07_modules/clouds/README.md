# Terraform Modules for multiple Clouds

In this example we have a Terraform Module to create a single VM on AWS, Azure or IBM Cloud. Depending of the selected module is where the VM is going to be created.

In Terraform it's not possible to dynamically select the source of a module and this makes sense, the module has to be loaded before executing the code and this load happens before knowing the value of the variable or condition that will define the module source.

This is different in Terranova. Based on the parameters you can load the code to execute into the platform. All the modules has to be accessible to the code at build time and at runtime the module is loaded.

- [Terraform Modules for multiple Clouds](#terraform-modules-for-multiple-clouds)
  - [Setup Access to the Clouds](#setup-access-to-the-clouds)
  - [Using the example](#using-the-example)
  - [Export the Terraform code](#export-the-terraform-code)

## Setup Access to the Clouds

Before the execution of the example you need an account in the cloud and provide the credentials or keys to access the cloud API.

It's required to have an account in the selected cloud environment, it can be a free account so you don't spend money executing the example. The creation and setup of an account is out of the scope of the document but the instructions are well documented in each cloud site.

The following links explain how to get and export the cloud credentials using a file.

- [AWS](./docs/access.md#aws)
- [Azure](./docs/access.md#azure)
- [IBM Cloud](./docs/access.md#ibm-cloud)

## Using the example

Use the `.credentials` script on every terminal/console that is going to apply the code so the provisioner can access to every cloud engine API (AWS, Azure or IBM Cloud).

## Export the Terraform code

Using the flag `--export` will generate and store the Terraform code into the directory `terraform_code` which you can use with Terraform.

Execute the following commands:

```bash
./terractl --export

source .credentials

cd terraform_code
terraform init
terraform plan
terraform apply

# On Azure, wait a few minutes to run:
terraform refresh

IP=$(terraform output public_ip)

ssh -i ~/.ssh/id_rsa ubuntu@$IP 'echo "Hello World"'

ssh -i ~/.ssh/id_rsa ubuntu@$IP
```

**Note**: On Azure, right after the first `terraform apply` the IP address is not displayed because it's not assigned yet, you have to refresh (`terraform refresh`) or apply the code again to get the assigned IP with `terraform output`.

To get the IP using the cloud CLI instead of Terraform (i.e. for debugging or to test the execution), use the line for the selected cloud.

```bash
# AWS:
IP=$()

# Azure:
IP=$(az vm show --resource-group terraformDemoResourceGroup --name terraformDemoVM -d --output tsv --query '[publicIps]')

# IBM Cloud:
IP=$()
```
