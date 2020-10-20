module github.com/johandry/terranova-examples/07_modules/clouds

go 1.15

require (
	// github.com/Azure/azure-sdk-for-go v36.2.0+incompatible // indirect
	// github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/johandry/terranova v0.0.4
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
// github.com/terraform-providers/terraform-provider-azurerm v1.34.0
)

replace github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281

replace github.com/johandry/terranova => ../../../terranova
