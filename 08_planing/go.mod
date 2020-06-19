module github.com/johandry/terranova-examples/02_ec2/ec2

go 1.13

require (
	github.com/hashicorp/terraform v0.12.20
	github.com/johandry/terranova v0.0.4
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
)

replace github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281

// use this replace when using or testing the local version of terranova
// replace github.com/johandry/terranova => ../../terranova
