# Terranova Examples

This repository contain examples about how to use the [Terranova](https://github.com/johandry/terranova) package. The examples are organized by complexity:

- **[Simple](./simple/)**: This is the simplest example, I recommend to start with it if this is the first time with Terranova. It is to create just one AWS EC2 instance. As simple as that. It can also be use it to test a local and modified version of Terranova.
- **[EC2](./ec2/)**: This example is a version of [Simple](./aws/simple/) but a little but complex. It create, scale or terminate AWS EC2 instances using the Terranova package. The build CLI receives parameters - such as the number of instances to create - making it more flexible and friendly to the user.
- **[Terraform](./terraform/)**: This example do not use Terranova, it use the Terraform package v0.11 just like Terranova uses it. Notice that it does not uses Terraform 0.12 neither the latest providers, use it for legacy Terraform code.
- **[MultiCloud](./multicloud/)**: This example create, scale or terminate a host either on AWS (EC2 instance), on VMWare (Virtual Machine) or OpenStack (Virtual Machine)
- **[Custom Logs](./custom-logs/)**: Terranova print to StdErr the verbosed logs from Terraform, in order to filter, modify or ignore such errors Terranova uses predefined Loggers or your own. This example has different Loggers in the [log/](./custom-logs/log/) directory that you can use as they are or use them to create your own Loggers.

All the examples have a README with all the instruction, also a Makefile to help you to execute all the possible actions with the example such as `build`, `test` and `run`, locally and on docker (i.e. `build-on-docker`). All of the examples build a binary and docker image named `terractl`, except the Simple example which binary name is `simple`.

