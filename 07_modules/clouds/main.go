package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/johandry/terranova"
	"github.com/johandry/terranova-examples/07_modules/clouds/data"
	"github.com/terraform-providers/terraform-provider-aws/aws"

	// "github.com/terraform-providers/terraform-provider-azurerm/azurerm"

	"github.com/johandry/terranova/logger"
)

const (
	exportDir     = "terraform_code"
	stateFilename = "cloud-vm.tfstate"
)

type options struct {
	Destroy    bool
	Export     bool
	Cloud      string
	PubKeyFile string
}

type output struct {
	IP string
}

func main() {
	opt, err := cliParse()
	if err != nil {
		die(err)
	}
	fmt.Printf("%+v\n", opt)

	out, err := apply(opt)
	if err != nil {
		die(err)
	}

	fmt.Printf("%+v\n", out)
}

func die(err error) {
	fmt.Printf("[ERROR] %s\n", err)
	os.Exit(1)
}

func cliParse() (*options, error) {
	opt := options{
		Cloud:      "aws",
		PubKeyFile: "~/.ssh/id_rsa.pub",
	}

	flag.BoolVar(&opt.Destroy, "destroy", opt.Destroy, "destroy the provisioned VM")
	flag.BoolVar(&opt.Export, "export", opt.Export, fmt.Sprintf("generate the Terraform code and export it to the %q directory", exportDir))
	flag.StringVar(&opt.Cloud, "cloud", opt.Cloud, "cloud where to provision the VM. The available options are: 'aws', 'az' and 'azure'")
	flag.StringVar(&opt.PubKeyFile, "pub", opt.PubKeyFile, "public key file to create the AWS Key Pair")
	flag.Parse()

	opt.Cloud = strings.ToLower(opt.Cloud)
	switch opt.Cloud {
	case "aws":
	case "az":
	case "azure":
		opt.Cloud = "az"
	default:
		return nil, fmt.Errorf("unsupported or unknown cloud %q", opt.Cloud)
	}

	return &opt, nil
}

func tfCode(opt *options) (string, error) {
	var tfCodeB bytes.Buffer
	tdTfCode := data.TfMainTmplData{
		Module: opt.Cloud,
	}
	tTfCode := template.Must(template.New("main.tf").Parse(data.TfMainTmpl))
	if err := tTfCode.Execute(&tfCodeB, tdTfCode); err != nil {
		return "", fmt.Errorf("fail to render the main.tf file. %s", err)
	}
	return tfCodeB.String(), nil
}

func apply(opt *options) (*output, error) {
	var logOutput io.Writer = os.Stderr
	myLog := logger.NewLog(logOutput, "WEB @ "+strings.ToUpper(opt.Cloud), logger.LogLevelInfo)
	logMiddleware := logger.NewMiddleware(myLog)
	defer logMiddleware.Close()

	tfMainCode, err := tfCode(opt)
	if err != nil {
		return nil, err
	}

	platform := terranova.NewPlatform(tfMainCode).
		Var("public_key_file", opt.PubKeyFile).
		SetMiddleware(logMiddleware)

	switch opt.Cloud {
	case "aws":
		platform.
			AddFile("modules/aws/main.tf", data.TfModAwsMainCode).
			AddProvider("aws", aws.Provider())
	case "az":
		platform.
			AddFile("modules/az/main.tf", data.TfModAzMainCode).
			AddProvider("aws", aws.Provider())
		// AddProvider("azurerm", azurerm.Provider())
	}

	if opt.Export {
		if err := os.MkdirAll(exportDir, 0755); err != nil {
			return nil, err
		}
		return nil, platform.Export(exportDir)
	}

	if _, err := platform.PersistStateToFile(stateFilename); err != nil {
		return nil, fmt.Errorf("fail to create the platform using state file %s. %s", stateFilename, err)
	}

	if err := platform.Apply(opt.Destroy); err != nil {
		return nil, fmt.Errorf("fail to apply the changes to the platform. %s", err)
	}

	if opt.Destroy {
		return nil, nil
	}

	ip, err := platform.OutputValueAsString("public_ip")
	if err != nil {
		return nil, err
	}

	out := &output{
		IP: ip,
	}

	return out, nil
}
