package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/johandry/terranova-examples/10_ibmcloud/letschat/data"
)

// TerraformTmplData is the template data structure
type TerraformTmplData struct {
	LetsChatPort string
	Status       bool
	StatusPort   string
	SSHAccess    bool
}

func tfCode(opt *options) (string, error) {
	var tfCodeB bytes.Buffer

	var terraformTmpl string
	switch opt.Cloud {
	case "aws":
		terraformTmpl = data.TerraformAWSTmpl
	case "ibm":
		terraformTmpl = data.TerraformIBMTmpl
	}
	tdTerraform := TerraformTmplData{
		LetsChatPort: opt.Port,
		Status:       opt.Status,
		StatusPort:   opt.StatusPort,
		SSHAccess:    opt.SSHAccess,
	}

	tTerraform := template.Must(template.New("main.tf").Parse(terraformTmpl))
	if err := tTerraform.Execute(&tfCodeB, tdTerraform); err != nil {
		return "", fmt.Errorf("failed render the terraform code. %v", err)
	}

	return tfCodeB.String(), nil
}
