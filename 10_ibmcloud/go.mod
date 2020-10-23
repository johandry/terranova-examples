module github.com/johandry/terranova-examples/10_ibmcloud/letschat

go 1.15

require (
	// github.com/IBM-Cloud/terraform-provider-ibm v1.13.1
	github.com/jarcoal/httpmock v1.0.6 // indirect
	github.com/johandry/terranova v0.0.4
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
)

replace github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281

replace github.ibm.com/ibmcloud/namespace-go-sdk => /Users/johandry/Workspace/sandbox/terraform-provider-ibm/common/github.ibm.com/ibmcloud/namespace-go-sdk

replace github.com/johandry/terranova => ../../terranova
