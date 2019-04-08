# Terranova Example: AWS EC2 instances

This example is used in the README of the Terranova repository and it is just to create one AWS EC2 instance.

Create the Key Pair with the following command:

```bash
aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
chmod 400 demo.pem
cat demo.pem
```

Login to the created EC2 instance using:

```bash
ssh -i ./demo.pem ubuntu@$(grep '"public_ip"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/')
```

To delete the Key Pair, execute:

```bash
aws ec2 delete-key-pair --key-name demo
rm demo.pem
```

To terminate the instance set the `count` variable to `0`, build it (`go build`) and execute it, or execute the following AWS CLI command:

```bash
id=$(grep '"id"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/' | uniq)
aws ec2 terminate-instances --instance-ids $id
```