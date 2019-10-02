# Terranova Example: AWS EC2 instances

This example is used in the README of the Terranova repository and it is just to create one AWS EC2 instance.

The `simple` application requires a Key Pair named `demo`. To create the Key Pair execute the following command:

```bash
aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
chmod 400 demo.pem
cat demo.pem
```

After the execution of the `simple` application you'll have an EC2 instance, the IP address can be obtained from the state file with:

```bash
grep '"public_ip"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/'
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

## Build

To build it just execute `go build .` and the output is the binary `simple` which when you execute it, you have an AWS EC2 instance running.

To build it in a isolated environment such as a Docker container, use `docker` and the `Dockerfile` to build the image `simple` like this:

```bash
docker build -t simple .
```

## Terranova local test

I use this simple example to do a quick test my local version of Terranova. To do this quick test I just use the `make` command:

```bash
make
```

The execution of `make` or `make build` get the local version of Terranova (located in `../../../terranova`), insert a `replace` in the `go.mod` file to use this local copy and build a Docker container with the binary of `simple`.

The Makefile also is useful to build and login into a Docker container with the environment to build the code (`make builder`) or login into a Docker container with the simple application (`make simple`).

To build locally and not in a Docker container, it's recommended to cleanup the modules before build. Then insert the `replace` in the `go.mod` file or use `make setup` to get the same:

```bash
go clean -modcache
make setup
go build .
make clean
```

