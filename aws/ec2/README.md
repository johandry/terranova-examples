# Terranova Example

This example is to create, scale or terminate AWS EC2 instances using the Terranova package. The code is explained in the blog post [Terranova: Using Terraform from Go](http://blog.johandry.com/post/terranova-terraform-from-go/).

To Build the example execute:

```bash
go build -o ec2 .
```

Before use the built binary `ec2` you need an AWS account, create it for free following [this instructions](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/). It's optional but suggested to install the AWS CLI application (instructions [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)).

It's required to have a Key Pair. To list all your available KeyPairs use the following AWS CLI command:

```bash
aws ec2 describe-key-pairs --query 'KeyPairs[*].KeyName' --output table
```

To create a new one, use this AWS CLI command:

```bash
aws ec2 create-key-pair --key-name MyKeyPair
```

To create or scale the created EC2 instances, use the command:

```bash
./ec2 --keyname MyKeyPair --count 3
```

This command will create 3 new EC2 instances if this is the first time running the command. If this is not the first time, it will scale up or down, creating or terminating existing instances to have a number of `3`.

To verify the existing EC2 instances, use the AWS CLI command:

```bash
aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table
```

To terminate the instances, use the command:

```bash
./ec2 --keyname MyKeyPair --count 0
```

Verify the results with the AWS CLI command above.

*IMPORTANT*: Do not delete the file `aws-ec2-ubuntu.tfstate` if you have EC2 instances. If you lost it, identify the instance IDs to terminate with the AWS CLI command previously used to verify results, then terminate your existing EC2 instances with the following AWS CLI command:

```bash
aws ec2 terminate-instances --instance-ids i-0c755d9c2d0d74f08 i-06519480197f8a82e
```

It's safe to delete the file `aws-ec2-ubuntu.tfstate` when the EC2 instances are terminated.