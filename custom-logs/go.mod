module github.com/johandry/terranova-examples/custom-logs

go 1.13

require (
	github.com/johandry/log v0.0.0-20190918193429-2b13006dd125
	github.com/johandry/terranova v0.0.3
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
)

replace github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281

// use this replace when using or testing the local version of terranova
// local: replace github.com/johandry/terranova => ../../terranova
// docker: replace github.com/johandry/terranova => ./pkg/terranova
