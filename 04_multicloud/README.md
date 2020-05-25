# Terranova Example: Multiple Clouds/Platforms

This example creates a number of VM on the requested platform or cloud, such as AWS, VMWare and OpenStack.

The value of the variable `platformName` which is set by the CLI parameter `--platform` will determine provider to setup (AWS, VMWare or OpenStack). Each platform has an init function to get the credentials from environment variables and use them as variables in the Terraform code to provision the required resources, such as VMs and key pairs.

The used providers and their version are in the `go.mod` file. Notice that the latest version of a provider may not work with Terranova, this is explained in the README of Terranova in the section **Providers version**.

## Build

To build it just execute `go build -o terractl .`, the output is the binary `terractl`.

## How to use

This example requires the credentials for the platforms to use in environment variables. Open the file `credentials.sh` and set the correct values for each platform. Then execute the file with `./credentials.sh`.

Execute the `terractl` binary using the following parameters:

* `--count` the number of instances to create or scale up/down. If it's `0` it will terminate all the existing instances.
* `--platform` is the platform to create, scale or destroy the VM or instances. The available platforms are: `aws`, `vsphere` and `openstack`. By default will use `aws`. 