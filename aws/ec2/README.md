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

*IMPORTANT*: Do not delete the file `aws-ec2-ubuntu.tfstate` if you have EC2 instances. If you lost it, identify the instance IDs to terminate with the AWS CLI command previously used to verify results, then terminate your existing EC2 instances and Key Pair with the following AWS CLI command:

```bash
aws ec2 terminate-instances --instance-ids i-0c755d9c2d0d74f08 i-06519480197f8a82e
aws ec2 delete-key-pair --key-name server_key
```

It's safe to delete the file `aws-ec2-ubuntu.tfstate` when the EC2 instances are terminated.