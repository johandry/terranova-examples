.PHONY: default test vet golint build init vendor

default: build

all: vendor build

aws-ec2-ubuntu: build

build: vet golint test
	@go build -o aws-ec2-ubuntu

test:
	@go list ./... | grep -v -E '^github.com/johandry/platformer-examples/aws/vendor' | xargs -n1 go test -cover

vet:
	@go list ./... | grep -v -E '^github.com/johandry/platformer-examples/aws/vendor' | xargs -n1 go vet -v

golint:
	-@[[ -x $${GOPATH}/bin/golint ]] || go get github.com/golang/lint/golint
	@go list ./... | grep -v -E '^github.com/johandry/platformer-examples/aws/vendor' | xargs -n1 golint

init:
	-@[[ -x $${GOPATH}/bin/govendor ]] || go get -u github.com/kardianos/govendor
	@govendor init
	@$(MAKE) vendor

vendor:
	@govendor list -no-status +missing | xargs -n1 go get -u
	@govendor add +external

update: vendor
	@govendor update +vendor
	@$(MAKE) build

clean:
	-@[[ ! -x ./vendor/vendor.json ]] || govendor remove +vendor
	@$(RM) -r aws-ec2-ubuntu

clean-all: clean
	@$(RM) -r ./vendor/
