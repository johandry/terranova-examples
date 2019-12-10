# Terranova Example: Custom Logs

This simple application requires a Key Pair named `demo`. To create the Key Pair execute the following command or `make key`:

```bash
aws ec2 create-key-pair --key-name demo --query 'KeyMaterial' --output text > demo.pem
chmod 400 demo.pem
cat demo.pem
```

After the execution of the application you'll have an EC2 instance, the IP address can be obtained from the state file with:

```bash
grep '"public_ip"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/'
```

Login to the created EC2 instance using `make ssh` or executing:

```bash
ssh -i ./demo.pem ubuntu@$(grep '"public_ip"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/')
```

To delete the Key Pair, execute `make key-clean` or:

```bash
aws ec2 delete-key-pair --key-name demo
rm demo.pem
```

To terminate the instance set the `count` variable to `0`, build it (`go build`) and execute it, or execute the following AWS CLI command:

```bash
id=$(grep '"id"' simple.tfstate | sed 's/.*: "\(.*\)",/\1/' | uniq)
aws ec2 terminate-instances --instance-ids $id
```

## How it works

This program is similar to the AWS/Simple example, the only difference is how to use the custom logs to print (or not) the intercepted logs from Terraform.

The value of the variable `logType` will define the type of custom log to use. The `switch` statement below will create the Terranova Log Middleware to intercept, parse and print the Terraform logs but using the created custom logger. 

The logger has to implement the Terranova interface `Logger` which defines certain number of methods such as `Infof`, `Warnf` and `Errors` among others, to print the log message in the desired format.

Each example custom logger is in a file located in the local package `log`, the different custom logs are explained below.

## Different custom logs

### Default

This is the default logger provided by Terranova. It just print every log entry in the format:

```
LEVEL [date] message
```

Where LEVEL is one of the following texts: `DEBUG`, `INFO`, `WARN` and `ERROR` . The date and message are the same from Terraform.

If the Middleware is created without any Logger or it is `nil`, the default logger will be used.

### Discard

This logger use the method `logger.DiscardLog()` to get a logger which will ignore every log output from Terraform. The entries will be intercepted and parsed but printed to the void.

Use this method when the Terraform output is not required.

### Terraform

This is just an example of what would happen when there is no Middleware defined, it just do nothing and prints every single terraform log output. It's noisy but you may want to do it this way to collect the information yourself or pass it as it's to the users.

To define the log level in Terraform, you can define the environment variable `TF_LOG` to one of the following values:  `TRACE`, `DEBUG`, `INFO`, `WARN` or `ERROR`. You will find this and more information at the [Debugging Terraform](https://www.terraform.io/docs/internals/debugging.html) page.

### Custom

You can define your own logger, just need to pass to the Terranova Log Middleware an interface of Logger, just like in the custom example.

This custom logger prints every event except for the Debug entries. You can also define how to print every entry, in this case, it's similar to the default logger.

Notice that the struct defines an instance of `log` and don't use the standard log instance. This is because the standard log instance is hijacked by the Middleware.

### `johandry/log` and log configuration

The package `github.com/johandry/log` allows you to create a custom logger from a viper (`github.com/spf13/viper`) object. Viper can provide configuration parameters from the Go code, environment variables, flags or program parameters and configuration files in different formats or source.

Here are 2 examples: `jlog-viper` configures the logger just from the code using a Viper instance, and `jlog-config` configures the logger from different sources: the Go code, environment variables and a YAML configuration file. Viper use the value based on the precedence defined in Viper and set the parameters with more priority.

Both examples use colors (configurable) to print the log entries and have the format:

```
[date] LEVEL Prefix: message
```

Where LEVEL is one of the following texts: `DEBUG`, `INFO`, `WARN` and `ERROR` . The date and message are the same from Terraform. The Prefix could be defined but by default it's `Terranova`.

The package `github.com/johandry/log` uses also `sirupsen/logrus` to format and print the log entries. 

### `sirupsen/logrus` in Text and JSON format

The `sirupsen/logrus` package is a widely used package for logging in Go. You can use it to print the logs in different formats, using fields, use Hooks and many other features. 

Here are 2 examples: `logrus-json` to print the log entries in JSON and `logrus` to print the log entries in text. Both examples print the fields `platform` and `count`, uses colors and begins with the log level in uppercase.

The JSON format may be interesting if another tool will read the output and process it.

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

The execution of `make` or `make build` get the local version of Terranova (located in `../../terranova`), insert a `replace` in the `go.mod` file to use this local copy and build a Docker container with the binary of `simple`.

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

## 