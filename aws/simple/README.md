# Terranova Example: AWS EC2 instances

This example is used in the README of the Terranova repository and it is just to create one AWS EC2 instance.

It uses the default logger, which will print only the Terraform info, warnning and errors log messages changing the format to: `LEVEL [date] message` where `LEVEL` could be `INFO`, `WARN` or `ERROR`.

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

To build it in a isolated environment such as a Docker container, use `docker` and the `Dockerfile` to build the container `simple` like this:

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

My workflow usually is:

Create the Key Pair:

```bash
make key
```

If building with Docker:

```bash
make
make run
make clean
```

If building locally:

```bash
make setup-local
make build-local
# Or
make run-simple
make clean-local
```

Test:

```bash
make list
make ssh
```

To cleanup, set the variable `count` to zero: `count := 0` and run the code. Then cleanup the rest:

```bash
make run           # or make run-local

make clean-key
make terminate     # Only if there are existing instances
rm simple simple.tfstate
```

## How it works

This is the most simplest way to use Terranova. The Terraform code is stored in the global variable `code` which is set in the `init()` method.

First we start defining a log instance because the use of the standard log will cause that any output is intercepted by the Terranova log middleware:

```go
log := log.New(os.Stderr, "", log.LstdFlags)
```

Then we create the log middleware to intercept every Terraform output, parse it and print it out using (in this case) the default and basic Terranova logger. This basic logger will print every Info, Warn and Error entry from Terranova.

To know more about the log middleware check the [EC2 example](../ec2) or the [custom-logs examples](../../custom-logs).

Starting this point, every log entry using the standard log is hijacked and parsed by the middleware. This hijack stops when we close the middleware:

```go
defer logMiddleware.Close()
```

The most important code segment is the one where the Platform is created:

```go
  platform, err := terranova.NewPlatform(code).
    AddMiddleware(logMiddleware).
    AddProvider("aws", aws.Provider()).
    Var("c", count).
    Var("key_name", keyName).
    ReadStateFromFile(stateFilename)
```

It starts passing the Terraform code and the next line sets the logger middleware. The following lines adds the AWS provider to the platform, because it's the provider used in the Terraform code. We send the value for the variables `c` and `key_name` which are also used in the code.

The last line is very important, here we load the platform state located in the file `simple.tfstate` (if exists). If there is on state file the requested amount of instances will be created, if there is a state then the instances will be scaled (up or down) or terminated (if the count is `0`).

At this point everything is set to apply the Terraform code, this is done with the method `Apply()` sending a boolean value to provision (`false`) or terminate (`true`).

Once the code is applied, we finalize with saving the current state into the state file, using the method `WriteStateToFile()`.
