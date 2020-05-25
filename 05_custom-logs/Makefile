# Local:
build:
	go build -o terractl .

run:
	go run .

setup:
	sed -i .bkp 's|^// replace github.com/johandry/terranova|replace github.com/johandry/terranova|' go.mod
	rm go.mod.bkp

clean-setup:
	sed -i .bkp 's|^replace github.com/johandry/terranova|// replace github.com/johandry/terranova|' go.mod
	rm -f go.mod.bkp

# On Docker:
build-on-docker:
	docker build -t terractl .

run-on-docker:
	docker run --rm -it terractl

# To Debug on Docker
sh-builder: setup
	docker build -t terractl-builder-dev --target=builder-dev .
	docker run --rm -it terractl-builder-dev

sh-app:
	docker build -t terractl-dev --target=app-dev .
	docker run --rm -it terractl-dev

# Login to EC2 instance
key:
	aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
	chmod 400 demo.pem
	cat demo.pem

ssh:
	ssh -i ./demo.pem ubuntu@$$(grep '"public_ip"' terractl.tfstate | sed 's/.*: "\(.*\)",/\1/')

# Clean Up
clean-key:
	aws ec2 delete-key-pair --key-name demo
	rm demo.pem

terminate:
	id=$(grep '"id"' terractl.tfstate | sed 's/.*: "\(.*\)",/\1/' | uniq)
	if [[ -n $$id ]]; then aws ec2 terminate-instances --instance-ids $$id; fi