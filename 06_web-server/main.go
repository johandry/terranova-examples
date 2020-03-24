package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/terraform/builtin/provisioners/file"
	"github.com/johandry/terranova"
	"github.com/johandry/terranova/logger"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var (
	tfCode      string
	port        string
	pubKeyFile  string
	privKeyFile string
	debug       bool
	quiet       bool
	destroy     bool
	export      bool
	sshAccess   bool
	sshFrom     string
)

const stateFilename = "aws-webserver.tfstate"

func main() {
	flag.Parse()

	log := log.New(os.Stderr, "", log.LstdFlags)

	myPublicIP := sshFrom
	if myPublicIP == "" {
		myPublicIP = getPublicIP()
		log.Printf("Using the public IP %s", myPublicIP)
	}

	if export {
		if err := ioutil.WriteFile("main.tf", []byte(tfCode), 0644); err != nil {
			log.Fatalf("failed to export the Terraform code to 'main.tf'. %v", err)
		}
		return
	}

	if debug && quiet {
		log.Fatal("debug mode and quiet mode cannot be set at the same time")
	}

	log.Printf("Starting (de)provisioning of the web server on AWS")

	logLevel := logger.LogLevelInfo
	if debug {
		logLevel = logger.LogLevelDebug
	}
	var output io.Writer = os.Stderr
	if quiet {
		output = ioutil.Discard
	}
	myLog := logger.NewLog(output, "WEB @ AWS", logLevel)
	logMiddleware := logger.NewMiddleware(myLog)
	defer logMiddleware.Close()

	platform, err := terranova.NewPlatform(tfCode).
		SetMiddleware(logMiddleware).
		AddProvider("aws", aws.Provider()).
		AddProvisioner("file", file.Provisioner()).
		Var("port", port).
		PersistStateToFile(stateFilename)

	if len(pubKeyFile) != 0 {
		platform.Var("public_key_file", pubKeyFile)
	}
	if len(privKeyFile) != 0 {
		platform.Var("private_key_file", privKeyFile)
	}

	if err != nil {
		log.Fatalf("Fail to create the platform using state file %s. %s", stateFilename, err)
	}

	if err := platform.Apply(destroy); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}

	if destroy {
		log.Printf("The web server has been destroyed")
		return
	}

	publicIP, _ := platform.OutputValueAsString("public_ip")
	log.Printf("Check the status and logs at: http://%s:%s", publicIP, "8081")
	log.Printf("The web server will be ready at: http://%s:%s", publicIP, port)

	if sshAccess {
		pubKeyFile, _ := platform.OutputValueAsString("public_key_file")
		log.Printf("Connect to the server executing: ssh -i %s %s@%s", pubKeyFile, "ubuntu", publicIP)
	}

}

func init() {
	flag.StringVar(&port, "port", "8080", "port to expose the webserver")
	flag.StringVar(&pubKeyFile, "pub", "", "public key file to create the AWS Key Pair")
	flag.StringVar(&privKeyFile, "priv", "", "private key file to connect to the new AWS EC2 instances")
	flag.BoolVar(&debug, "debug", false, "debug mode, prints also debug output from terraform")
	flag.BoolVar(&quiet, "quiet", false, "quiet/silence mode, do not print any terraform output")
	flag.BoolVar(&destroy, "destroy", false, "terminate the web server instance(s)")
	flag.BoolVar(&export, "export", false, "export the Terraform code to the file main.tf")
	flag.BoolVar(&sshAccess, "shh-access", true, "enable SSH access to the hosts")
	flag.StringVar(&sshFrom, "from", "", "Allow connection from this IP address. If empty, the public IP address will be used")

	tfCode = `
variable "port" { default = 8080 }
variable "logport" { default = 8081 }
variable "public_key_file"  { default = "~/.ssh/id_rsa.pub" }
variable "private_key_file" { default = "~/.ssh/id_rsa" }
locals {
	public_key    = "${file(pathexpand(var.public_key_file))}"
	private_key   = "${file(pathexpand(var.private_key_file))}"
	userdata   		= <<-USERDATA
#! /bin/bash

# Start the status reporter
html_status() { 
	echo '<li><span style="color:blue">[INFO]</span> '$1'</li>' >> status.txt; 
	echo "<html><head><title>Server Status</title>" > index.html;
	[[ $2 != "END" ]] && echo '<meta http-equiv="refresh" content="10"/>' >> index.html; 
	echo '</head><body><h2>Status:</h2><ul style="list-style-type:none">' >> index.html; 
	cat status.txt >>index.html;
	echo '</ul><h2>Logs:</h2><textarea rows="50" cols="100">' >>index.html
	cat /var/log/cloud-init-output.log >>index.html
	echo '</textarea>' >>index.html
	echo "</body></html>" >>index.html; 
}
nohup busybox httpd -f -p ${var.logport} &
echo $! > httpd.pid

# Install Docker
html_status "updating the packages..."
sudo apt-get update
html_status "done"
html_status "installing curl, CA certificates and other packages..."
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
html_status "done"
html_status "installing docker..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo usermod -aG docker $USER
html_status "done"

# Install Docker Compose
html_status "installing docker compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
html_status "done"

# Create the Docker Compose file
html_status "creating docker compose file..."
cat <<DOCKERFILE >docker-compose.yaml
version: '3'
services:
  app:
    image: sdelements/lets-chat:latest
    links:
      - mongo
    ports:
      - ${var.port}:8080
      - 5222:5222

  mongo:
    image: mongo:latest
DOCKERFILE
html_status "done"

html_status "starting docker compose with Let'sChat and MongoDB..."
# kill -9 $(cat httpd.pid)
# rm -f httpd.pid

# Start Docker Compose
sudo docker-compose up -d
html_status "done"

TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/public-hostname)

html_status 'Go to: <a href="http://'$IP':${var.port}">http://'$IP':${var.port}</a>' "END"
USERDATA
}

output "public_ip" {
  value   = "${aws_instance.webserver.public_ip}" 
}
output "public_key_file" {
  value   = "${var.private_key_file}" 
}

provider "aws" {
  region   = "us-west-2"
}

resource "aws_instance" "webserver" {
  ami                    = "ami-0d1cd67c26f5fca19"
	instance_type          = "t2.micro"
	key_name      				 = "webserver_key"
	vpc_security_group_ids = [ "${aws_security_group.webserver_ingress.id}" ]
	user_data_base64       = base64encode(local.userdata)

  tags = {
    Name = "terranova-example-webserver"
	}
}

resource "aws_key_pair" "keypair" {
	key_name    = "webserver_key"
	public_key  = local.public_key
}

resource "aws_security_group" "webserver_ingress" {
  name = "terranova_webserver_ingress"

  ingress {
		description = "Open port to web application hosted"
    from_port   = var.port
    to_port     = var.port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
	}

  ingress {
		description = "Open port to web application hosted"
    from_port   = var.logport
    to_port     = var.logport
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
	}

  ingress {
		description = "Open port to SSH access"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
	}
	
	egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
`
}

// Get the public IP address using one of the following API services:
// https://www.ipify.org
// http://myexternalip.com
// http://api.ident.me
// http://whatismyipaddress.com/api
func getPublicIP() string {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(ip)
}
