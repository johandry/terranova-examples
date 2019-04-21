# Terranova Example: Using Terraform package directly

In this example is used the Terraform package just like Terranova uses it. It's more work and many other things are not included.

Use the `Makefile` to build, use or test the application.

## Build

To build just need to execute `make build`, optionally set the environment variable `APP_NAME` if you want another name but `terraform`.

To recreate the modules, execute `make mod`, this will delete the current module files, the cache, recreate the module files and download them all.

## Use

After build the code, you'll see the binary `terraform`, use it to create, scale or destroy EC2 instances in AWS. It's required to have an AWS account as well as the AWS credentials exposed with the AWS environment variables or with AWS CLI credentials configured.

Flags:

* `count`: number of instances to create or scale. If set to `0` will terminate all the created instances. Default is `2`.
* `pub`: SSH public key used to create the AWS Key Pair. Default is `~/.ssh/id_rsa.pub`.
* `priv`: SSH Private key file to connect to the new AWS EC2 instances. Default is `~/.ssh/id_rsa`.
* `state`: Location of the Terraform state file. This is required to scale or terminate instances. Default is `./tf.state`.

Example:

```bash
terraform --count 3
terraform --count 1
terraform --count 0
```

## Validations

List the existing Key Pairs and EC2 instances with `make describe-key-pairs describe-instances` or the following commands:

```bash
aws ec2 describe-key-pairs --query 'KeyPairs[*].KeyName' --output table
aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table
```

Login to the created EC2 instance using `make ssh IP=<ip address>` or the following command and assuming you already have your SSH keys created:

```bash
ssh -i ~/.ssh/id_rsa -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ubuntu@$(IP)
```

Or, you may run some automated tests with `make test` but this assumes all your EC2 instances were created with this application.

## Terminate

As mentioned before, execute `terraform --count 0` to terminate all the created instances. If for some reason this process fail, you may terminate manually all the resources `make terminate-instances delete-key-pair` or with the following commands:

```bash
aws ec2 terminate-instances --instance-ids $$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --output text | tr '\n' ' ')
aws ec2 delete-key-pair --key-name server_key
```

These commands assumes all your instances where create with this application.

<!-- 
## In Docker

Use the `Makefile` to do some build, use and test the application in an isolated environment such as a Docker container.

Use `make docker-build` to build the `terraform` image locally. To use it, execute the command `make docker-run` and this command requires the AWS credentials exposed with environment variables or having the AWS CLI credentials configured. Optionally, use the variable `COUNT` to specify the number of instances to create, scale or terminate.

```bash
export AWS_ACCESS_KEY_ID=<your AWS access key>
export AWS_SECRET_ACCESS_KEY=<your AWS secret access key>

make docker-run COUNT=3
make docker-run COUNT=0
``` 
-->