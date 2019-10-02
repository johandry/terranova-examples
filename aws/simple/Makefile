
build: setup
	docker build -t simple .
	# Use: docker run --rm -it simple

builder: setup
	docker build -t simple-builder-dev --target=builder-dev .
	docker run --rm -it simple-builder-dev

simple: setup
	docker build -t simple-dev --target=app-dev .
	docker run --rm -it simple-dev

setup:
	sed -i .bkp 's|^// replace github.com/johandry/terranova|replace github.com/johandry/terranova|' go.mod
	rm go.mod.bkp
	rm -rf ./pkg/
	mkdir ./pkg
	cp -r ../../../terranova ./pkg/

clean:
	rm -rf ./pkg
	sed -i .bkp 's|^replace github.com/johandry/terranova|// replace github.com/johandry/terranova|' go.mod
	rm -f go.mod.bkp
	-docker image rm simple simple-builder-dev simple-dev
