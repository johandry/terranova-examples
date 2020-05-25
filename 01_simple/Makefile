# For Build Local:
build:
	go build .

run:
	go run .

# To use local package
setup:
	sed -i .bkp 's|^// local: replace github.com/johandry/terranova|replace github.com/johandry/terranova|' go.mod
	rm go.mod.bkp

clean-setup:
	sed -i .bkp 's|^replace github.com/johandry/terranova|// local: replace github.com/johandry/terranova|' go.mod
	rm -f go.mod.bkp

# To Build on Docker:
build-on-docker: setup
	docker build -t simple .

run-on-docker:
	docker run --rm -it simple

# To use local package on Docker
setup-for-docker:
	sed -i .bkp 's|^// docker: replace github.com/johandry/terranova|replace github.com/johandry/terranova|' go.mod
	rm go.mod.bkp
	rm -rf ./pkg/
	mkdir ./pkg
	cp -r ../../../terranova ./pkg/

clean-setup-for-docker:
	rm -rf ./pkg
	sed -i .bkp 's|^replace github.com/johandry/terranova|// docker: replace github.com/johandry/terranova|' go.mod
	rm -f go.mod.bkp
	-docker image rm simple simple-builder-dev simple-dev

# To Debug on Docker
sh-builder: setup
	docker build -t simple-builder-dev --target=builder-dev .
	docker run --rm -it simple-builder-dev

sh-app: setup
	docker build -t simple-dev --target=app-dev .
	docker run --rm -it simple-dev

# To use the application 
key:
	aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
	chmod 400 demo.pem
	cat demo.pem

clean-key:
	aws ec2 delete-key-pair --key-name demo
	rm demo.pem

terminate:
	id=$(grep '"id"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/' | uniq)
	if [[ -n $$id ]]; then aws ec2 terminate-instances --instance-ids $$id; fi

# Checking/using the instance(s)
ssh:
	ssh -i ./demo.pem ubuntu@$$(grep '"public_ip"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/')

list:
	aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table

set-one:
	sed -i.bkp 's/count := 0/count := 1/' main.go

set-zero:
	sed -i.bkp 's/count := 1/count := 0/' main.go

clean:
	rm simple simple.tfstate
	rm main.go.bkp

all: key set-one setup build run list ssh set-zero build run set-one clean-key clean-setup clean

all-on-docker: key set-one setup-for-docker build-on-docker run-on-docker list ssh set-zero build-on-docker run-on-docker set-one clean-key clean-setup-for-docker clean
