module github.com/johandry/terranova-examples/aws/simple

go 1.13

require (
	github.com/johandry/terranova v0.0.0-20190422213246-704ed6ce88e7
	github.com/terraform-providers/terraform-provider-aws v1.60.0
)

// use this replace when using or testing the local version of terranova
replace github.com/johandry/terranova => ./pkg/terranova
