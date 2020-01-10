package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform/builtin/provisioners/file"
	"github.com/hashicorp/terraform/terraform"
	"github.com/johandry/terranova"
	"github.com/terraform-providers/terraform-provider-aws/aws"
	"github.com/terraform-providers/terraform-provider-openstack/openstack"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere"
)

var (
	code         string
	vars         map[string]interface{}
	platformName string
	count        int
	pubKeyFile   string
	privKeyFile  string
)

func main() {
	flag.Parse()

	if count < 0 {
		log.Fatalf("count cannot be negative. It has to be either '0' to terminate all the created instances or the desired number of instances")
	}

	var provider terraform.ResourceProvider

	switch platformName {
	case "aws":
		provider = aws.Provider()
		initAWS()
	case "vsphere":
		provider = vsphere.Provider()
		initVsphere()
	case "openstack":
		provider = openstack.Provider()
		initOpenStack()
	default:
		log.Fatalf("unknown platform name %q, use one of the following platforms: \"aws\", \"vsphere\", \"openstack\"", platformName)
	}

	stateFilename := platformName + ".tfstate"

	platform, err := terranova.NewPlatform(code).
		AddProvider(platformName, provider).
		AddProvisioner("file", file.Provisioner()).
		Var("srv_count", count).
		BindVars(vars).
		PersistStateToFile(stateFilename)

	if err != nil {
		log.Fatalf("Fail to create the platform using state file %s. %s", stateFilename, err)
	}

	if err := platform.Apply((count == 0)); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}

	log.Printf("Your virtual machine instances is ready")
}

func init() {
	flag.IntVar(&count, "count", 2, "number of instances to create. Set to '0' to terminate them all.")
	flag.StringVar(&platformName, "platform", "aws", "cloud or platform to create the instance")
}
