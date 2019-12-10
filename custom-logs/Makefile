# To Build on Docker:

build: setup
	docker build -t simple .
	# Use: docker run --rm -it simple

run:
	docker run --rm -it simple

builder: setup
	docker build -t simple-builder-dev --target=builder-dev .
	docker run --rm -it simple-builder-dev

simple: setup
	docker build -t simple-dev --target=app-dev .
	docker run --rm -it simple-dev

setup:
	sed -i .bkp 's|^// docker: replace github.com/johandry/terranova|replace github.com/johandry/terranova|' go.mod
	rm go.mod.bkp
	rm -rf ./pkg/
	mkdir ./pkg
	cp -r ../../terranova ./pkg/

clean:
	rm -rf ./pkg
	sed -i .bkp 's|^replace github.com/johandry/terranova|// docker: replace github.com/johandry/terranova|' go.mod
	rm -f go.mod.bkp
	-docker image rm simple simple-builder-dev simple-dev

# For Build Local:
build-local:
	go build .

setup-local:
	sed -i .bkp 's|^// local: replace github.com/johandry/terranova|replace github.com/johandry/terranova|' go.mod
	rm go.mod.bkp

clean-local:
	sed -i .bkp 's|^replace github.com/johandry/terranova|// local: replace github.com/johandry/terranova|' go.mod
	rm -f go.mod.bkp

# To Run:
run-simple:
	go run .

key:
	aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
	chmod 400 demo.pem
	cat demo.pem

ssh:
	ssh -i ./demo.pem ubuntu@$$(grep '"public_ip"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/')

clean-key:
	aws ec2 delete-key-pair --key-name demo
	rm demo.pem

terminate:
	id=$(grep '"id"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/' | uniq)
	if [[ -n $$id ]]; then aws ec2 terminate-instances --instance-ids $$id; fi