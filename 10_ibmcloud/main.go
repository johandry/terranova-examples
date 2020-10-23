package main

import (
	"log"
)

func main() {
	opt := cliParse()
	tfCode, err := tfCode(opt)
	if err != nil {
		log.Fatalf(err.Error())
	}

	switch {
	case opt.Export:
		log.Printf("Starting export of Terraform code to %s", exportDir)
	case opt.Destroy:
		log.Printf("Starting deprovisioning of the LetsChat server on IBM Cloud")
	default:
		log.Printf("Starting provisioning of the LetsChat server on IBM Cloud")
	}

	publicIP, err := apply(tfCode, opt)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case opt.Export:
		log.Printf("The Terraform code has been exported to %s", exportDir)
	case opt.Destroy:
		log.Printf("The LetsChat server has been destroyed")
	default:
		log.Printf("The LetsChat server will be ready at: http://%s:%s", publicIP, opt.Port)
		if opt.Status {
			log.Printf("Check the status and logs at: http://%s:%s", publicIP, opt.StatusPort)
		}
		if opt.SSHAccess {
			log.Printf("Connect to the server executing: ssh -i %s %s@%s", opt.PrivKeyFile, "ubuntu", publicIP)
		}
	}
}
