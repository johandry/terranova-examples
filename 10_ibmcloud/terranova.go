package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	// "github.com/IBM-Cloud/terraform-provider-ibm/ibm"
	"github.com/johandry/terranova"
	"github.com/johandry/terranova-examples/10_ibmcloud/letschat/data"
	"github.com/johandry/terranova/logger"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

const (
	exportDir     = "./terraform_code"
	stateFilename = "webserver.tfstate"
)

func apply(tfCode string, opt *options) (string, error) {
	logMiddleware := getLogMiddleware(opt)
	defer logMiddleware.Close()

	platform := terranova.NewPlatform(tfCode).
		AddFile("user_data.sh", data.UserdataTmpl).
		AddFile("docker-compose.yaml", data.DockerComposeTmpl).
		SetMiddleware(logMiddleware)

	switch opt.Cloud {
	case "aws":
		platform.AddProvider("aws", aws.Provider())
	case "ibm":
		// platform.AddProvider("ibm", ibm.Provider())
	}

	if opt.Status {
		platform.Var("public_key_file", opt.PubKeyFile)
	}

	if opt.Export {
		if err := os.MkdirAll(exportDir, 0755); err != nil {
			return "", err
		}
		return "", platform.Export(exportDir)
	}

	if _, err := platform.PersistStateToFile(stateFilename); err != nil {
		return "", fmt.Errorf("Fail to create the platform using state file %s. %s", stateFilename, err)
	}

	if err := platform.Apply(opt.Destroy); err != nil {
		return "", fmt.Errorf("Fail to apply the changes to the platform. %s", err)
	}

	if opt.Destroy {
		return "", nil
	}

	return platform.OutputValueAsString("public_ip")
}

func getLogMiddleware(opt *options) *logger.Middleware {
	logLevel := logger.LogLevelInfo
	if opt.Debug {
		logLevel = logger.LogLevelDebug
	}
	var output io.Writer = os.Stderr
	if opt.Quiet {
		output = ioutil.Discard
	}

	myLog := logger.NewLog(output, "WEB @ IBM Cloud", logLevel)

	return logger.NewMiddleware(myLog)
}
