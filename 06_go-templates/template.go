package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"text/template"

	"github.com/johandry/terranova-examples/06_go-templates/letschat/data"
)

func tfCode(opt *options) (string, error) {
	userdataB64, err := userData(opt)
	if err != nil {
		return "", err
	}

	var tfCodeB bytes.Buffer
	tdTerraform := data.TerraformTmplData{
		UserdataB64:  userdataB64,
		LetsChatPort: opt.Port,
		StatusPort:   opt.StatusPort,
		SSHAccess:    opt.SSHAccess,
	}

	tTerraform := template.Must(template.New("main.tf").Parse(data.TerraformTmpl))
	if err := tTerraform.Execute(&tfCodeB, tdTerraform); err != nil {
		return "", fmt.Errorf("failed render the terraform code. %v", err)
	}

	return tfCodeB.String(), nil
}

func userData(opt *options) (string, error) {
	dockerComposeB64, err := dockerCompose(opt)
	if err != nil {
		return "", err
	}

	var userdataB bytes.Buffer
	tdUserdata := data.UserdataTmplData{
		LetsChatPort:     opt.Port,
		StatusPort:       opt.StatusPort,
		DockerComposeB64: dockerComposeB64,
	}

	tUserdata := template.Must(template.New("userdata.sh").Parse(data.UserdataTmpl))
	if err := tUserdata.Execute(&userdataB, tdUserdata); err != nil {
		return "", fmt.Errorf("failed render the userdata code. %v", err)
	}

	return base64.StdEncoding.EncodeToString(userdataB.Bytes()), nil
}

func dockerCompose(opt *options) (string, error) {
	var dockerComposeB bytes.Buffer
	tdDockerCompose := data.DockerComposeTmplData{
		LetsChatPort: opt.Port,
	}

	tDockerCompose := template.Must(template.New("docker-compose.yaml").Parse(data.DockerComposeTmpl))
	if err := tDockerCompose.Execute(&dockerComposeB, tdDockerCompose); err != nil {
		return "", fmt.Errorf("failed render the Docker Compose yaml. %v", err)
	}

	return base64.StdEncoding.EncodeToString(dockerComposeB.Bytes()), nil
}
