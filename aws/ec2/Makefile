KEYNAME 	?= terranova

build:
	go build -o ec2 .

clean:
	$(RM) ec2

clean-all: clean
	$(RM) *.tfstate

describe-key-pairs:
	aws ec2 describe-key-pairs --query 'KeyPairs[*].KeyName' --output table

create-key-pair:
	aws ec2 create-key-pair --key-name $(KEYNAME)

describe-instances:
	aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table

terminate-instances:
	aws ec2 terminate-instances --instance-ids $$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --output text | tr '\n' ' ')

	