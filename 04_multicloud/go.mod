module github.com/johandry/terranova-examples/04_multicloud/vm

go 1.15

require (
	github.com/hashicorp/terraform v0.12.20
	github.com/johandry/terranova v0.0.4
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
	github.com/terraform-providers/terraform-provider-openstack v1.23.0
	github.com/terraform-providers/terraform-provider-vsphere v1.13.0
)

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.2.0
	github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281
)

// use this replace when using or testing the local version of terranova
// replace github.com/johandry/terranova => ../../../terranova
