module github.com/johandry/terranova-examples/01_simple/simple

go 1.15

require (
	github.com/johandry/terranova v0.0.4
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
)

replace github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281

// use this replace when using or testing the local version of terranova
// local: replace github.com/johandry/terranova => ../../terranova

// docker: replace github.com/johandry/terranova => ./pkg/terranova
