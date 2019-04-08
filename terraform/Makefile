SHELL			:= /bin/bash

build:
	go build -o tf .

clean:
	$(RM) tf

clean-all: clean
	$(RM) *.tfstate

describe-key-pairs:
	aws ec2 describe-key-pairs --query 'KeyPairs[*].KeyName' --output table

delete-key-pair:
	aws ec2 delete-key-pair --key-name server_key

describe-instances:
	aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table

terminate-instances:
	aws ec2 terminate-instances --instance-ids $$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --output text | tr '\n' ' ')

ssh-keys:
	ssh-keygen -t rsa

ssh:
	ssh -i ~/.ssh/id_rsa -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ubuntu@$(IP)

test:
	@expected="ami used: ami-6e1a0117"; \
	result="OK"; \
	for ip in $$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].PublicIpAddress' --output text | tr '\n' ' '); do \
		echo -n "Checking $$ip: "; \
		actual=$$(ssh -q -i ~/.ssh/id_rsa -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ubuntu@$$ip cat /tmp/file.log); \
		if [[ $${actual} == $${expected} ]]; then echo "OK"; else echo "FAIL"; result="FAIL"; fi; \
	done; echo; \
	echo "Test results: $${result}"