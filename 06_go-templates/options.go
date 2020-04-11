package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type options struct {
	Debug       bool
	Quiet       bool
	Destroy     bool
	Export      bool
	Port        string
	Status      bool
	StatusPort  string
	PubKeyFile  string
	PrivKeyFile string
	SSHAccess   bool
	SSHFrom     string
}

func cliParse() *options {
	myPublicIP := getPublicIP()
	opt := options{
		Port:        "8080",
		StatusPort:  "8081",
		PubKeyFile:  "~/.ssh/id_rsa.pub",
		PrivKeyFile: "~/.ssh/id_rsa",
		SSHFrom:     myPublicIP,
	}

	flag.BoolVar(&opt.Debug, "debug", opt.Debug, "debug mode, prints also debug output from terraform")
	flag.BoolVar(&opt.Quiet, "quiet", opt.Quiet, "quiet/silence mode, do not print any terraform output")
	flag.BoolVar(&opt.Destroy, "destroy", opt.Destroy, "terminate the web server instance(s)")
	flag.BoolVar(&opt.Export, "export", opt.Export, "export the Terraform code to the file main.tf")
	flag.StringVar(&opt.Port, "port", opt.Port, "port to expose the chat application")
	flag.BoolVar(&opt.Status, "status", opt.Status, "shows the server status during the setup")
	flag.StringVar(&opt.StatusPort, "status-port", opt.StatusPort, "port to expose the server status")
	flag.StringVar(&opt.PubKeyFile, "pub", opt.PubKeyFile, "public key file to create the AWS Key Pair")
	flag.StringVar(&opt.PrivKeyFile, "priv", opt.PrivKeyFile, "private key file to connect to the new AWS EC2 instances")
	flag.BoolVar(&opt.SSHAccess, "ssh-access", opt.SSHAccess, "enable SSH access to the hosts")
	flag.StringVar(&opt.SSHFrom, "ssh-from", opt.SSHFrom, "Allow connection from this IP address")

	flag.Parse()

	// If --status-port is set then --status is implicitly set
	for _, f := range os.Args[1:] {
		if strings.TrimLeft(f, "-") == "status-port" {
			opt.Status = true
			break
		}
	}

	if opt.Debug && opt.Quiet {
		log.Fatal("debug mode and quiet mode cannot be set at the same time")
	}

	return &opt
}

// Use one of the following API services to get the public IP address:
// https://www.ipify.org?format=text
// http://myexternalip.com/raw
// https://v4.ident.me/
// http://ipv4bot.whatismyipaddress.com
const myIPURL = "https://www.ipify.org?format=text"

// Get the public IP address
func getPublicIP() string {
	resp, err := http.Get(myIPURL)
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
