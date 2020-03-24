SHELL			:= /bin/bash

APP_NAME		?= tf
IMG_NAME		?= tform

build:
	go build -o $(APP_NAME) .

clean:
	$(RM) $(APP_NAME)

clean-all: clean
	$(RM) *.tfstate

mod: mod-clean
	go mod init github.com/johandry/terranova-examples/terraform
	go mod tidy

mod-clean:
	$(RM) go.mod go.sum
	go clean -modcache

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

docker-build:
	docker build -t $(APP_NAME) .

COUNT									?= 1
AWS_ACCESS_KEY_ID 		?= $(shell cat ~/.aws/credentials | grep 'default]' -A 2 | grep aws_access_key_id | cut -d= -f2 | tr -d ' ')
AWS_SECRET_ACCESS_KEY ?= $(shell cat ~/.aws/credentials | grep 'default]' -A 2 | grep aws_secret_access_key | cut -d= -f2 | tr -d ' ')

docker-run:
	mkdir -p tmp
	docker run \
		-v $(HOME)/.ssh:/root/terraform/ssh \
		-v $(PWD)/tmp:/root/terraform \
		-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
		-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
		$(IMG_NAME) \
			--state /root/terraform/tf.state \
			--pub /root/terraform/ssh/id_rsa.pub \
			--priv /root/terraform/ssh/id_rsa \
			--count $(COUNT)

docker-bash:
	mkdir -p tmp
	docker run -it \
		-v $(HOME)/.ssh:/root/terraform/ssh \
		-v $(PWD)/tmp:/root/terraform \
		-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
		-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
		--entrypoint "/bin/sh" \
		$(IMG_NAME)