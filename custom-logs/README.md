# Terranova Example: Custom Logs

This example is similar to the Simple example, the only difference is how to use the custom logs to print (or not) the intercepted logs from Terraform.

The value of the variable `logType` will define the type of custom log to use. The `switch` statement below will create the Terranova Log Middleware to intercept, parse and print the Terraform logs but using the created custom logger.

The logger has to implement the Terranova interface `Logger` which defines certain number of methods such as `Infof`, `Warnf` and `Errors` among others, to print the log message in the desired format.

Each example custom logger is in a file located in the local package `log`, the different custom logs are explained below.

## Different custom logs

### Default

This is the default logger provided by Terranova. It just print every log entry in the format:

```text
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

```text
[date] LEVEL Prefix: message
```

Where LEVEL is one of the following texts: `DEBUG`, `INFO`, `WARN` and `ERROR` . The date and message are the same from Terraform. The Prefix could be defined but by default it's `Terranova`.

The package `github.com/johandry/log` uses also `sirupsen/logrus` to format and print the log entries.

### `sirupsen/logrus` in Text and JSON format

The `sirupsen/logrus` package is a widely used package for logging in Go. You can use it to print the logs in different formats, using fields, use Hooks and many other features.

Here are 2 examples: `logrus-json` to print the log entries in JSON and `logrus` to print the log entries in text. Both examples print the fields `platform` and `count`, uses colors and begins with the log level in uppercase.

The JSON format may be interesting if another tool will read the output and process it.

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

## 