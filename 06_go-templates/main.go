package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	opt := cliParse()
	tfCode, err := tfCode(opt)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if opt.Export {
		if err := ioutil.WriteFile("main.tf", []byte(tfCode), 0644); err != nil {
			log.Fatalf("failed to export the Terraform code to 'main.tf'. %v", err)
		}
		return
	}

	log := log.New(os.Stderr, "", log.LstdFlags)
	log.Printf("Starting (de)provisioning of the LetsChat server on AWS")

	publicIP, err := apply(tfCode, opt)
	if err != nil {
		log.Fatal(err)
	}

	if opt.Destroy {
		log.Printf("The web server has been destroyed")
		return
	}

	if opt.Status {
		log.Printf("Check the status and logs at: http://%s:%s", publicIP, opt.StatusPort)
	}

	log.Printf("The web server will be ready at: http://%s:%s", publicIP, opt.Port)

	if opt.SSHAccess {
		log.Printf("Connect to the server executing: ssh -i %s %s@%s", opt.PrivKeyFile, "ubuntu", publicIP)
	}
}
