# Terranova Example: Getting started

This example is used in the README of the Terranova repository and it is just to create one AWS EC2 instance.

It uses the default logger, which will print only the Terraform info, warnning and errors log messages changing the format to: `LEVEL [date] message` where `LEVEL` could be `INFO`, `WARN` or `ERROR`.

## Build

Before build, modify the following line to set the logger to use in the example

```go
logType := "custom"
```

The available options are: `default`, `terraform`, `discard`, `custom`, `jlog-viper`, `jlog-config`, `logrus`, `logrus-json`.

Optionally you can modify the value of `count` to set the number of instances to create. Setting this variable to `0` will terminate the created instances.

```go
count := 1
```

To build it just execute `go build -o terractl .` or `make`, the output is the binary `terractl`.

To build it in a isolated environment such as a Docker container, you'll need `docker`, then execute `make build-on-docker` or:

```bash
docker build -t terractl .
```

## How to use

This example requires a Key Pair named `demo`. To create the Key Pair execute the following command or `make key`:

```bash
aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
chmod 400 demo.pem
cat demo.pem
```

Execute the example using one of the following lines:

```bash
./terractl
# Or:
go run .
# Or:
docker run --rm -it terractl
```

Also you can use `make run` and `make run-ondocker`.

After the execution of the application you'll have an EC2 instance, the IP address can be obtained from the state file with:

```bash
grep '"public_ip"' terractl.tfstate | sed 's/.*: "\(.*\)",/\1/'
```

Login to the created EC2 instance using `make ssh` or executing:

```bash
ssh -i ./demo.pem ubuntu@$(grep '"public_ip"' terractl.tfstate | sed 's/.*: "\(.*\)",/\1/')
```

To delete the Key Pair, execute `make key-clean` or:

```bash
aws ec2 delete-key-pair --key-name demo
rm demo.pem
```

To terminate the instance set the `count` variable to `0`, and execute:

```bash
go run .
```

You may terminate the instance using the AWS CLI:

```bash
id=$(grep '"id"' terractl.tfstate | sed 's/.*: "\(.*\)",/\1/' | uniq)
aws ec2 terminate-instances --instance-ids $id
rm *.tfstate
```

## Terranova development and test

This simple example can be used to test the local version of the Terranova package. To do this test with Docker execute the following commands:

```bash
make setup-for-docker
make build-on-docker
make run-on-docker

# Cleanup
make clean-setup-for-docker
```

The same can be done without docker executing the following commands:

```bash
make setup
make build
make run

# Cleanup
make clean-setup
```

The execution of `make setup` or `make setup-for-docker` gets the local version of Terranova (located in `../../terranova`), insert a `replace` in the `go.mod` file to use this local copy and build a Docker container with the binary of `simple`.

The Makefile also is useful to build and login into a Docker container with the environment to build the code (`make sh-builder`) or login into a Docker container with the simple application (`make sh-app`).

To build locally and not in a Docker container, it's recommended to cleanup the modules before build.

```bash
go clean -modcache
make setup
go build .
make clean-setup
```

The common workflow would be something like this:

```bash
make key
sed -i.bkp 's/count := 0/count := 1/' main.go
make setup  			# Or: setup-for-docker

# Build & run
make build				# Or: build-on-docker
make run					# Or: run-on-docker

# Test:
make list
make ssh

# Terminate
sed -i.bkp 's/count := 1/count := 0/' main.go
make build				# Or: build-on-docker
make run					# Or: run-on-docker

# If the execution fails, execute:
make terminate

# Cleanup
sed -i.bkp 's/count := 0/count := 1/' main.go
make clean-key
rm simple simple.tfstate
rm main.go.bkp

make clean-setup	# Or: clean-setup-for-docker
```

Or execute `make all` or `make all-on-docker`.

## How it works

This is the most simplest way to use Terranova. The Terraform code is stored in the global variable `code` which is set in the `init()` method.

First we start defining a log instance because the use of the standard log will cause that any output is intercepted by the Terranova log middleware:

```go
log := log.New(os.Stderr, "", log.LstdFlags)
```

Then we create the log middleware to intercept every Terraform output, parse it and print it out using (in this case) the default and basic Terranova logger. This basic logger will print every Info, Warn and Error entry from Terranova. This hijack stops when we close the middleware:

```go
defer logMiddleware.Close()
```

To know more about the log middleware check the [EC2 example](../ec2) or the [custom-logs examples](../../custom-logs).

The most important code segment is when the Platform is created:

```go
  platform, err := terranova.NewPlatform(code).
    SetMiddleware(logMiddleware).
    AddProvider("aws", aws.Provider()).
    Var("c", count).
    Var("key_name", keyName).
    PersistStateToFile(stateFilename)
```

It starts passing the Terraform code and the next line sets the logger middleware. The following lines adds the AWS provider to the platform, because it's the provider used in the Terraform code. We send the value for the variables `c` and `key_name` which are also used in the code.

The last line is very important, here we load the platform state located in the file `simple.tfstate` (if exists). If there is no state file the requested amount of instances will be created, if there is a state then the instances will be scaled (up or down) or terminated (if the count is `0`). The function `PersistStateToFile()` ensures the state file is always updated when a change is done.

At this point everything is set to apply the Terraform code, this is done with the method `Apply()` sending a boolean value to provision (`false`) or terminate (`true`).
