# Setup Access to the Clouds

- [Setup Access to the Clouds](#setup-access-to-the-clouds)
  - [AWS](#aws)
  - [Azure](#azure)
  - [IBM Cloud](#ibm-cloud)

## AWS

To start using AWS you need the *AWS Access Key ID* and *AWS Secret Access Key*, follow these instructions to get them:

1. Sign in to the AWS Management Console and open the IAM console at https://console.aws.amazon.com/iam/.
2. In the navigation pane, choose **Users**.
3. Choose the name of the user whose access keys you want to create, and then choose the **Security credentials** tab.
4. In the **Access keys** section, choose **Create access key**.
5. Download the key pair choosing **Download .csv file**. Store the keys in a secure location. You will not have access to the secret access key again after this dialog box closes.

To setup the AWS access is very simple once you have installed the AWS CLI tool. Execute `aws configure` then enter the *AWS Access Key ID* and *AWS Secret Access Key* as they are in the downloaded `.csv` file. Select the default region closer to where you are and the default output format, they can be `json`, `yaml`, `text` and `table`, choose the one you are comfortable using.

This is enough to execute the example but you can store the credentials in the `.credetials` files with all the other cloud credentials, like so.

```bash
#!/bin/sh

echo "Setting AWS environment variables ..."
export AWS_ACCESS_KEY_ID=<AWS Access Key ID>
export AWS_SECRET_ACCESS_KEY=<AWS Secret Access Key>
```

## Azure

Before using the Azure module you need to login to Azure. Execute the following commands:

```bash
az login
az account list --query "[].{name:name, subscriptionId:id, tenantId:tenantId}"
```

Get the value of `subscriptionId` to export it in the variable `SUBSCRIPTION_ID`, then execute the following command to create a role assignment:

```bash
export SUBSCRIPTION_ID='00000000-0000-0000-0000-000000000000'
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/${SUBSCRIPTION_ID}"
```

The output of the last command look like this:

```json
{
  "appId": "00000000-0000-0000-0000-000000000000",
  "displayName": "azure-cli-2020-04-29-00-11-22",
  "name": "http://azure-cli-2020-04-29-00-11-22",
  "password": "00000000-0000-0000-0000-000000000000",
  "tenant": "00000000-0000-0000-0000-000000000000"
}
```

Use this output to create the file `.credentials` with the following content:

```bash
#!/bin/sh

echo "Setting Azure environment variables ..."
export ARM_SUBSCRIPTION_ID=<subscriptionId>
export ARM_CLIENT_ID=<appId>
export ARM_CLIENT_SECRET=<password>
export ARM_TENANT_ID=<tenant>

export ARM_ENVIRONMENT=public
```

To verify Terraform can access your Azure account execute the following Terraform code, executing the Terraform commands: `terraform init && terraform plan && terraform apply`

```hcl
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "rg" {
  name     = "testResourceGroup"
  location = "westus"
}
```

## IBM Cloud

There are many resources you can create for free on IBM Cloud but these examples requires to upgrade the account and pay for some resources. So, we suggest to run the IBM Cloud example if you have an upgraded account.

It's recommended to install the IBM Cloud CLI and required plugins to have a second option besides Terraform to access the cloud, on macOS and Linux that can be done with the execution of the following commands.

```bash
curl -sL https://ibm.biz/idt-installer | bash

ibmcloud version        # to verify it's successfully installed
ibmcloud plugin install -f -r "IBM Cloud" infrastructure-service
```

Execute the following commands to login it to your account.

```bash
ibmcloud login --sso
ibmcloud target         # to verify it's successfully logged in
```

With Terraform 0.12, the IBM Cloud Terraform provider has to be installed manually. Download the latest version from https://github.com/IBM-Cloud/terraform-provider-ibm/releases for your system and unzip it into the directory `~/.terraform.d/plugins/`.

_Note: On macOS you may get an error the first time trying to use the provider, to solve this error open a Finder window on `~/.terraform.d/plugins/` and execute the provider from there Finder right-clicking on the file, then selecting **Open with ...** > **Terminal**, click on **Open** in the new pop-up window. After closing the Terminal window where the provider ran, it can be used without any problem._

```bash
unzip terraform-provider-ibm_*.zip
mv terraform-provider-ibm_v* ~/.terraform.d/plugins/
open ~/.terraform.d/plugins/
```

Create an API Key to be used by Terraform and store it in the environment variable `IC_API_KEY`. This can be done on the Web Console or with the following commands. Once you generate the API Key it's not possible to view it again, so safe the file `.terraform_key.json` with the API Key to a secure location.

```bash
ibmcloud iam api-key-create TerraformKey -d "API Key for Terraform $(date +%m-%d-%Y)" --file .terraform_key.json
export IC_API_KEY=$(grep '"apikey":' .terraform_key.json | sed 's/.*: "\(.*\)".*/\1/')
```

Using the IBM Cloud web console follow these steps:

1. Go to **Manage** > **Access (IAM)** > **API keys**.
2. Click **Create an IBM Cloud API key**.
3. Enter a name and description for your API key, then click **Create**.
4. Click on the show icon to display the API key. Or, click **Copy** to copy and save it for later, or click **Download** to have a file like in the example above.

For more information read the [Creating an API key](https://cloud.ibm.com/docs/account?topic=account-userapikey#create_user_key) documentation.

You may export the variable `IC_API_KEY` in your shell profile file (i.e. `~/.bashrc` or `~/.zshrc`), or export it using the same `.credentials` file, like this.

```bash
#!/bin/sh

echo "Setting IBM Cloud environment variables ..."
export IC_API_KEY=<IBM Cloud API Key>
```

To verify Terraform can access your IBM Cloud account execute the following Terraform code, executing the Terraform commands: `terraform init && terraform plan && terraform apply`

```hcl
provider "ibm" {
  generation         = 2
  region             = "us-south"
}

data "ibm_iam_access_group" "accgroup" {
}

output "accgroups" {
  value = data.ibm_iam_access_group.accgroup.groups[*].name
}
```
