module github.com/johandry/terranova-examples/aws/ec2

go 1.13

require (
	github.com/hashicorp/terraform v0.12.12
	github.com/johandry/terranova v0.0.2
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
)

// use this replace when using or testing the local version of terranova
// replace github.com/johandry/terranova => ../../../terranova
