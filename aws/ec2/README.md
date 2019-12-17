# Terranova Example: AWS EC2 instances

This example is to create, scale or terminate AWS EC2 instances using the Terranova package. The code is explained in the blog post [Terranova: Using Terraform from Go](http://blog.johandry.com/post/terranova-terraform-from-go/).

To Build the example execute:

```bash
go build -o ec2 .
```

Before use the built binary `ec2` you need an AWS account, create it for free following [this instructions](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/). It's optional but suggested to install the AWS CLI application (instructions [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)).

It's required to have a private and public keys. If you are on Mac OSX execute `ssh-keygen -t rsa`. By default the private and public keys are located in `~/.ssh/id_rsa` and `~/.ssh/id_rsa.pub`.

To create or scale the created EC2 instances, use the command:

```bash
./ec2 --count 3
```

This command will create 3 new EC2 instances if this is the first time running the command. If this is not the first time, it will scale up or down, creating or terminating existing instances to have a number of `3`. The program is also assuming the default private and public keys, if they are not the defaults, provide them with the flags: `--pub /path/to/public/key.pem --priv /path/to/private/key.pem`.

To verify the existing EC2 instances, use the AWS CLI command:

```bash
aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table
```

To login to an instance, use the `ssh` command as follows:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@<IP Address>
```

Get the IP address with the verification command above and make sure to specify the correct path for the private key with the `-i` flag. If you cannot login, verify you have an inbound rule in the default security group to allow connection to port 22 from anywhere (not recommended) or your network address.

To terminate the instances, use the command:

```bash
./ec2 --count 0
```

Verify the results with the AWS CLI command above.

The binary accept the following parameters:

* `--count` the number of instances to create or scale up/down. If it's `0` it will terminate all the existing instances.
* `--pub` and `--priv` are the public or private key files. If they are passed as parameters they won't be created.
* `--debug` prints the Terraform debug logs. If not used (default) the binary only prints info, warning and error log level entries from Terraform.
* `--quiet` do not print any Terraform output, it will print only the output from the binary.

*IMPORTANT*: Do not delete the file `aws-ec2-ubuntu.tfstate` if you have EC2 instances. If you lost it, identify the instance IDs to terminate with the AWS CLI command previously used to verify results, then terminate your existing EC2 instances and Key Pair with the following AWS CLI command:

```bash
aws ec2 terminate-instances --instance-ids i-0c755d9c2d0d74f08 i-06519480197f8a82e
aws ec2 delete-key-pair --key-name server_key
```

It's safe to delete the file `aws-ec2-ubuntu.tfstate` when the EC2 instances are terminated.

## How it works

The program starts with the setup and parse of the flags, then it sets the global variable `code` with the Terraform code to execute. The Terraform code creates a key pair used to login into the created instances, then create the amount of requested instances and finally creates the file `/tmp/file.log` in every instance with the AMI Id used.

After the parse of the flags (line `flag.Parse()`), a log instance is created. This instance is necessary if we are going to use the log middleware.

The log middleware hijacks the standard log instance to intercept the Terraform output, parse it and send it back to the custom logger (`myLog`). This custom logger uses the default Terranova logger sending the output to StdErr or discard it if the flag `--quiet` was used. The custom logger also uses the log level Info (prints Info, Warns and Errors) or the log level Debug (prints debug entries + the same as Info) if the flag `--debug` is set.

Starting this point, every log entry using the standard log is hijacked and parsed by the middleware. This hijack stops when we close the middleware:

```go
defer logMiddleware.Close()
```

One of the most important code segment is the one where the Platform is created:

```go
  platform, err := terranova.NewPlatform(code).
    SetMiddleware(logMiddleware).
    AddProvider("aws", aws.Provider()).
    AddProvisioner("file", file.Provisioner()).
    Var("srv_count", count).
    ReadStateFromFile(stateFilename)
```

It starts passing the Terraform code and the next line sets the logger middleware. The following lines adds the AWS provider to the platform, because it's the provider used in the Terraform code. The same with the File provisioner, because it's used in the Terraform code (`provisioner "file"`) and send the value for the variable `srv_count` which is also used in the code (`variable "srv_count"`).

The last line is very important, here we load the platform state located in the file `aws-ec2-ubuntu.tfstate` (if exists). If there is on state (file) the requested amount of instances will be created, if there is a state then the instances will be scaled (up or down) or terminated (if the amount is `0`).

If a variable have a default value in the Terraform code, it's optional to set its value using the `Var()` method. That's what happens with the variables `public_key_file` and `private_key_file`.

At this point everything is set to apply the Terraform code, this is done with the method `Apply()` sending a boolean value to provision (`false`) or terminate (`true`).

Once the code is applied, we finalize with saving the current state into the state file, using the method `WriteStateToFile()`.
