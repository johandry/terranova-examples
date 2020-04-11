# Terranova Example: Using Go Templates

This example is a bit complex compare to the others however more practical. This example will create a free tier instance on AWS hosting a [Let'sChat](https://sdelements.github.io/lets-chat/) service. Let'sChat is a self-hosted chat for small teams, this example is also ephemeral, once the instance/service is destroyed there is no trace of the chats - in theory.

This example uses Go templates to generate the Terraform code, this is a feature that is not - yet - in Terraform. The use of Go templates makes the Terraform variables optional and we can do amazing things such as repeat code using for-loops or insert code if some conditions are accomplished.

The structure of the example is as follow:

```text
.
├── apply.go
├── data
│   ├── docker_compose_yaml.go
│   ├── main_tf.go
│   └── userdata_sh.go
├── go.mod
├── main.go
├── options.go
└── template.go
```

The `main.go` file use the code in `options.go` to parse all the CLI parameters and return a struct with all the options provided by the user or from the default values. After that uses the code in `template.go` to build the Terraform code. The code is created from the templates located in the files in the `data/` directory, each `.go` file contain a struct with the template variables and the template itself.

The Terraform code (located in `main_tf.go`) create an EC2 instance in a security group inside the default VPC (I know, it's insecure but I'm trying to make this example simple). The security group allow ingress traffic through the port where Let'sChat is running and egress traffic everywhere.

If the CLI option `--shh-access` is set, the user can login to the instance with SSH. If so, the Terraform code creates Key Pair from the private/public keys `~/.ssh/id_rsa` and `~/.ssh/id_rsa.pub` or the keys passed in the parameters `--priv` and `--pub`. Notice that the definition of the variables for the private and public keys depends if the SSH access is requested or not. Also, if the SSH access is allowed, the security group will have an ingress rule to allow access to port 22 from the local public IP or the IP given in the parameter `--ssh-from`.

The startup of the Let'sChat application takes some time, in order to know the status of this task the logs are exposed in the port `8081` or the port given in the parameter `--status-port`. This is an insecure option, that port may be used for malicious uses but it may be useful for debugging or the first time using the application and you don't know what is happening. The status report is disable by default, to enable it the  `--status` parameter is required. If so, the Terraform code will include an ingress rule in the security group to allow traffic thought the port where the status is exposed (defaulted to `8081`).

The use of Go templates may reduce the Terraform code and therefore remove responsibilities to Terraform, this might help us to have a faster execution of Terraform and give us more control of what the program is doing.

Once the Terraform code is build, it can be exported, saving the `main.tf` file with all the code or it could be executed, using the function `apply()` in the file `apply.go`.

The export of the code can be used for debugging, if you modify the code and it isn't working, you can export it and try to execute it with `terraform` so you can identify what is failing: the Terraform code built or the Go code with Terranova.

There is no import feature in this application and I think it does not makes any sense. If there is a Terraform code in files, then why would you need Terranova? Instead use `terraform` and you got the results. However, the Terraform code may be in the GitHub repository and with the use of Go generators that code can be embedded into the Go code, just like we have it now. So, instead of modifying the Terraform code in the Go code, we modify the Terraform code/files then execute `go generate` to build the Go code. There will be an example of this later.

The execution is just like the previous simple examples. It creates the log middleware to intercept the Terraform logs and print them in a different way, creates the platform with the built code, the AWS provider and sending the Terraform state to a file. If the status report is required, the variables with the path to the private and public key are set. After the code is applied we get the Terraform output using `platform.OutputValueAsString()` and return it.

The `main.go` finalize printing the URL to access the Let'sChat application or reporting that the instance is destroyed (if `--destroy` was set). It may also display the URL to view the status or the `ssh` command to SSH into the EC2 instance if any of these were required.

## Build & Usage

To build just execute `go build .`, the built binary name is `letschat`.

The basic execution `letschat` then `letschat --destroy` when the chat application is no needed. Optionally you can use the parameter `--status` to view the startup logs, this process takes some time so it's recommended the first time you use it; the parameter `--shh-access` to access with `ssh` to the instance for debugging and see what's happening, or `--export` to view or study the Terraform code generated. To view the other parameters use `--help`.
