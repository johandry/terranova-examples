package data

// TfModAzMainCode is the plain text code for module/az/main.tf
const TfModAzMainCode = `
variable "public_key" {}

locals {
  location = "eastus"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "terraform_demo_group" {
  name     = "terraformDemoResourceGroup"
  location = local.location

  tags = {
    environment = "Terraform Demo"
  }
}

resource "azurerm_virtual_network" "terraform_demo_network" {
  name                = "terraformDemoVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.terraform_demo_group.location
  resource_group_name = azurerm_resource_group.terraform_demo_group.name

  tags = {
    environment = "Terraform Demo"
  }
}

resource "azurerm_subnet" "terraform_demo_subnet" {
  name                 = "terraformDemoSubnet"
  resource_group_name  = azurerm_resource_group.terraform_demo_group.name
  virtual_network_name = azurerm_virtual_network.terraform_demo_network.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "terraform_demo_public_ip" {
  name                = "terraformDemoPublicIP"
  location            = azurerm_resource_group.terraform_demo_group.location
  
  resource_group_name = azurerm_resource_group.terraform_demo_group.name
  allocation_method   = "Dynamic"

  tags = {
    environment = "Terraform Demo"
  }
}

resource "azurerm_network_security_group" "terraform_demo_nsg" {
  name                = "terraformDemoNetworkSecurityGroup"
  location            = azurerm_resource_group.terraform_demo_group.location
  resource_group_name = azurerm_resource_group.terraform_demo_group.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "Terraform Demo"
  }
}

resource "azurerm_network_interface" "terraform_demo_nic" {
  name                = "terraformDemoNIC"
  location            = azurerm_resource_group.terraform_demo_group.location
  resource_group_name = azurerm_resource_group.terraform_demo_group.name

  ip_configuration {
    name                          = "terraformDemoNicConfiguration"
    subnet_id                     = azurerm_subnet.terraform_demo_subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.terraform_demo_public_ip.id
  }

  tags = {
    environment = "Terraform Demo"
  }
}

# Connect the security group to the network interface
resource "azurerm_network_interface_security_group_association" "terraform_demo_nsg_nic" {
  network_interface_id      = azurerm_network_interface.terraform_demo_nic.id
  network_security_group_id = azurerm_network_security_group.terraform_demo_nsg.id
}

resource "random_id" "randomId" {
  keepers = {
    resource_group = azurerm_resource_group.terraform_demo_group.name
  }
  byte_length = 8
}

resource "azurerm_storage_account" "terraformDemoStorageAccount" {
  name                     = "diag${random_id.randomId.hex}"
  resource_group_name      = azurerm_resource_group.terraform_demo_group.name
  location                 = azurerm_resource_group.terraform_demo_group.location
  account_replication_type = "LRS"
  account_tier             = "Standard"

  tags = {
    environment = "Terraform Demo"
  }
}

resource "azurerm_linux_virtual_machine" "terraform_demo_vm" {
  name                  = "terraformDemoVM"
  location              = azurerm_resource_group.terraform_demo_group.location
  resource_group_name   = azurerm_resource_group.terraform_demo_group.name
  network_interface_ids = [azurerm_network_interface.terraform_demo_nic.id]
  size                  = "Standard_DS1_v2"

  os_disk {
    name                 = "terraformDemoOsDisk"
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  computer_name                   = "terraformdemovm"
  admin_username                  = "ubuntu"
  disable_password_authentication = true

  admin_ssh_key {
    username   = "ubuntu"
    public_key = var.public_key
  }

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.terraformDemoStorageAccount.primary_blob_endpoint
  }

  tags = {
    environment = "Terraform Demo"
  }
}

output "public_ip" {
  value = azurerm_public_ip.terraform_demo_public_ip.ip_address
}
`
