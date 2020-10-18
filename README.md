# Terranova Examples

This repository contain examples about how to use the [Terranova](https://github.com/johandry/terranova) package. The examples are organized by complexity:

1. **[Simple](./01_simple/)**: This is the simplest example, I recommend to start with it if this is the first time with Terranova. It is to create just one AWS EC2 instance. As simple as that. It can also be use it to test a local and modified version of Terranova.
2. **[EC2](./02_ec2/)**: This example is a version of [Simple](./aws/simple/) but a little but complex. It create, scale or terminate AWS EC2 instances using the Terranova package. The build CLI receives parameters - such as the number of instances to create - making it more flexible and friendly to the user.
3. **[Terraform](./03_terraform/)**: This example do not use Terranova, it use the Terraform package v0.11 just like Terranova uses it. Notice that it does not uses Terraform 0.12 neither the latest providers, use it for legacy Terraform code.
4. **[MultiCloud](./04_multicloud/)**: This example create, scale or terminate a host either on AWS (EC2 instance), on VMWare (Virtual Machine) or OpenStack (Virtual Machine)
5. **[Custom Logs](./05_custom-logs/)**: Terranova print to StdErr the verbosed logs from Terraform, in order to filter, modify or ignore such errors Terranova uses predefined Loggers or your own. This example has different Loggers in the [log/](./custom-logs/log/) directory that you can use as they are or use them to create your own Loggers.
6. **[Web Server](./06_web-server/)**: Creates a Web Server instance running and hosting the given web content.
7. **[Terraform Modules](./07_modules)**: There are a few examples using Terraform Modules with Terranova.
   1. **[Cloud Providers](./07_modules/clouds)**: With Terraform it's not possible to select dynamically a module, with Terranova you can choose the module to load based on the user input. All modules has the same input and output parameters so the use is standard the only change is the source to know which module to load.
   2. **[Environments](./07_modules/environments)**: This is a simple example to deploy a web application into different environments (production and staging). The application is the same but the environment requirements are different, depending of the user input is the environment used to deploy the web application.

All the examples have a README with all the instruction, also a Makefile to help you to execute all the possible actions with the example such as `build`, `test` and `run`, locally and on docker (i.e. `build-on-docker`). All of the examples build a binary and docker image named `terractl`, except the Simple example which binary name is `simple`.
