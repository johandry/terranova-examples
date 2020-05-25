package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/johandry/terranova-examples/06_go-templates/letschat/data"
)

func tfCode(opt *options) (string, error) {
	var tfCodeB bytes.Buffer
	tdTerraform := data.TerraformTmplData{
		LetsChatPort: opt.Port,
		Status:       opt.Status,
		StatusPort:   opt.StatusPort,
		SSHAccess:    opt.SSHAccess,
	}

	tTerraform := template.Must(template.New("main.tf").Parse(data.TerraformTmpl))
	if err := tTerraform.Execute(&tfCodeB, tdTerraform); err != nil {
		return "", fmt.Errorf("failed render the terraform code. %v", err)
	}

	return tfCodeB.String(), nil
}
