package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/builtin/provisioners/file"
	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

// Code is the Terraform code to execute
const Code = `
variable "count" 	          { default = 2 }
variable "public_key_file"  { default = "~/.ssh/id_rsa.pub" }
variable "private_key_file" { default = "~/.ssh/id_rsa" }
locals {
  public_key    = "${file(pathexpand(var.public_key_file))}"
  private_key   = "${file(pathexpand(var.private_key_file))}"
}
provider "aws" {
  region        = "us-west-2"
}
resource "aws_instance" "server" {
  instance_type = "t2.micro"
  ami           = "ami-6e1a0117"
  count         = "${var.count}"
  key_name      = "server_key"

  provisioner "file" {
    content     = "ami used: ${self.ami}"
    destination = "/tmp/file.log"

    connection {
      user        = "ubuntu"
      private_key = "${local.private_key}"
    }
  }
}
resource "aws_key_pair" "keypair" {
  key_name    = "server_key"
  public_key  = "${local.public_key}"
}
`

const (
	stateFile = "tf.state"
)

var (
	count       int
	pubKeyFile  string
	privKeyFile string
)

// Platformer store all the information needed by Terraform
type Platformer struct {
	Code         string
	Vars         map[string]interface{}
	Providers    map[string]terraform.ResourceProvider
	Provisioners map[string]terraform.ResourceProvisioner
	State        *terraform.State
}

func main() {
	flag.IntVar(&count, "count", 2, "number of instances to create. Set to '0' to terminate them all.")
	flag.StringVar(&pubKeyFile, "pub", "", "public key file to create the AWS Key Pair")
	flag.StringVar(&privKeyFile, "priv", "", "private key file to connect to the new AWS EC2 instances")
	flag.Parse()

	var state bytes.Buffer

	p := &Platformer{
		Code: Code,
		Vars: map[string]interface{}{
			"count": count,
		},
	}

	if len(pubKeyFile) != 0 {
		p.Vars["public_key_file"] = pubKeyFile
	}
	if len(privKeyFile) != 0 {
		p.Vars["private_key_file"] = privKeyFile
	}

	// If the file exists, read the state from the state file
	if _, errStat := os.Stat(stateFile); errStat == nil {
		stateB, err := ioutil.ReadFile(stateFile)
		if err != nil {
			log.Fatalf("Fail to read the state from %q", stateFile)
		}
		state = *bytes.NewBuffer(stateB)

		// Get the Terraform state from the state file content
		if p.State, err = terraform.ReadState(&state); err != nil {
			log.Fatalln(err)
		}
	}

	// Create a temporal directory or use any directory
	tfDir, err := ioutil.TempDir("", ".tf")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(tfDir)
	// Save the code into a single or multimple files
	filename := filepath.Join(tfDir, "main.tf")
	configFile, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer configFile.Close()

	// Copy the Terraform template from p.Code into the new file
	if _, err = io.Copy(configFile, strings.NewReader(p.Code)); err != nil {
		log.Fatalln(err)
	}

	// Create the Terraform module
	mod, err := module.NewTreeModule("", tfDir)
	if err != nil {
		log.Fatalln(err)
	}

	// Create the Storage pointing to where the Terraform code is
	storageDir := filepath.Join(tfDir, "modules")
	s := module.NewStorage(storageDir, nil)
	s.Mode = module.GetModeNone // or module.GetModeGet

	// Finally make the module load the
	if err := mod.Load(s); err != nil {
		log.Fatalf("Failed loading modules. %s", err)
	}

	// Optionally, you can validate the loaded code if it has some user input
	if err := mod.Validate().Err(); err != nil {
		log.Fatalf("Failed Terraform code validation. %s", err)
	}

	// Add Providers:
	ctxProviders := make(map[string]terraform.ResourceProviderFactory)
	// ctxProviders["null"] = terraform.ResourceProviderFactoryFixed(null.Provider())
	ctxProviders["aws"] = terraform.ResourceProviderFactoryFixed(aws.Provider())
	providerResolvers := terraform.ResourceProviderResolverFixed(ctxProviders)

	// Add Provisioners:
	provisionersFactory := make(map[string]terraform.ResourceProvisionerFactory)
	provisionersFactory["file"] = func() (terraform.ResourceProvisioner, error) {
		return file.Provisioner(), nil
	}

	destroy := (count == 0)

	ctxOpts := terraform.ContextOpts{
		Destroy:          destroy,
		State:            p.State,
		Variables:        p.Vars,
		Module:           mod,
		ProviderResolver: providerResolvers,
		Provisioners:     provisionersFactory,
	}

	ctx, err := terraform.NewContext(&ctxOpts)
	if err != nil {
		log.Fatalf("Failed creating context. %s", err)
	}

	if _, err := ctx.Refresh(); err != nil {
		log.Fatalln(err)
	}
	if _, err := ctx.Plan(); err != nil {
		log.Fatalln(err)
	}
	if _, err := ctx.Apply(); err != nil {
		log.Fatalln(err)
	}

	// Retrive the state from the Terraform context
	p.State = ctx.State()
	if err := terraform.WriteState(p.State, &state); err != nil {
		log.Fatalf("Failed to retrive the state. %s", err)
	}
	// Save the state to the local file 'tf.state'
	if err = ioutil.WriteFile(stateFile, state.Bytes(), 0644); err != nil {
		log.Fatalf("Fail to save the state to %q. %s", stateFile, err)
	}

}
