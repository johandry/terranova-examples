# Terranova: Using Terraform from Go

**Terranova** = **Terraform** + **Go**

![](https://github.com/johandry/terranova/raw/master/images/terranova-160x160.png)

## What is Go?

https://golang.org

![](https://golang.org/lib/godoc/images/go-logo-blue.svg)



```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
}
```

Version: 1.15

Packages: https://golang.org/pkg/

Playground: https://play.golang.org

Main features:

* Concurrency
* Simplicity and Consistency
* Binaries
* Powerful standard library
* Package Management

## What is Terraform?

https://www.terraform.io/

![](https://www.terraform.io/assets/images/logo-hashicorp-3f10732f.svg)

*"Terraform is a tool for building, changing, and versioning infrastructure"*

The **Infrastructure Code** is written in a custom domain-specific-language (DSL) called **HCL** (*Hashicorp Configuration Language*), it's similar to JSON and YAML. With HCL we write declarative definitions for the infrastructure we want to exist. Then the Terraform engine takes care of provisioning and updating resources.

```hcl
provider "ibm" {
  generation         = 2
  region             = var.region
}

variable "region" {
  default = "us-south"
}

data "ibm_iam_access_group" "accgroup" {
}

output "accgroups" {
  value = data.ibm_iam_access_group.accgroup.groups[*].name
}
```

With HCL we define **what** we want and Terraform make sure the infrastructure is as we want it to be. It's like Kubernetes and the manifest YAML files.

To use Terraform we need:

1. Terraform
2. The Terraform Provider, in this example is IBM Cloud Provider
3. A cloud account and access to the cloud API, for IBM Cloud this is the API Key for VPC Gen 2

Having the code, we just execute:

```bash
> terraform init
Initializing the backend...

Initializing provider plugins...
  ...
Terraform has been successfully initialized!
  ...
> terraform plan
Refreshing Terraform state in-memory prior to plan...
  ...

data.ibm_iam_access_group.accgroup: Refreshing state...

  ...
❯ terraform apply
data.ibm_iam_access_group.accgroup: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

accgroups = [
  "Add cases and view orders",
  ...
  "View account summary",
  "View cases",
]
```

In this example we did not build or provision any resource, this code just query the API for the IAM Access Groups and print the results.

To know more about Terraform, Infrastructure as Code on IBM Cloud, read: https://ibm.github.io/cloud-enterprise-examples/iac/content-overview and for the simplified examples, go to https://github.com/IBM/cloud-enterprise-examples/tree/master/iac 

## What is Terranova?

https://github.com/johandry/terranova

**Terranova** = **Terraform** + **Go**

Terranova is a Go package that allows you to easily use the Terraform Go Packages instead of executing the Terraform binary

![](https://github.com/johandry/terranova/raw/master/images/terranova-160x160.png)

Terranova is a Go package to use Terraform in your Go code

* Terraform with HCL => Declarative definitions       => WHAT

* Terranova with Go   => Imperative programming  => HOW

### Benefits

1. **Release and deliver a single binary**: Instead of a bunch of Terraform infrastructure code there is only one binary to execute.
2. **No dependencies**: with a single binary the user don't need to install Terraform, the Providers neither the infrastructure code.
3. **Less errors**: the user cannot modify the infrastructure code directly reducing the errors due to user misconfiguration or incorrect input values
4. **Go programming language**: HCL is getting more features and improvements making possible to developers to add more logical sentences like loops, conditionals but will never be comparable to a programming language like Go.
5. **Better integration with external components**:  You can do whatever is possible with Go or any programming language: 
   1. Configuration files, in multiple formats: YAML, JSON, TOML, Properties, ...
   2. Expose or consult an API, using REST, GRPC, GraphQL or any API method supported in Go.
   3. Integrate your application with other tools such as Consul, Vault, ...
   4. Use of external libraries or Go packages which functionalities are not provided by HCL, Terraform or the providers.
6. **Better User Interface**: Terraform is mostly CLI, you can have Terraform Cloud or IBM Schematics but they are not customizable. Using your own Go program it could be a CLI, Web, Application, Mobile, ...
7. **Templates and dynamic infrastructure code**: Terraform supports templates with external files, you can use templates with Terraform or Go Templates. The infrastructure code can be a template rendered by the application then pass it to Terraform.
8. **Mix Declarative & Imperative programming**: You still use HCL to declarative define the infrastructure totally or partially, and use Go to define how to execute some processes. You get the benefits of both methods.
9. **Terraform**: It's still Terraform! Terranova is just a Go package using the Terraform code, you are still using Terraform, but not the `terraform` binary.

### Alternative options

* **Cloud SDK for Go**: Using the [AWS SDK](https://aws.amazon.com/tools/) or [IBM Cloud SDK](https://cloud.ibm.com/docs?tab=api-docs) is an option but the code is 100% imperative, there is no HCL. This gives you more control of the code and it may be harder to maintain. If your program is multi-cloud you may need to import multiple SDK libraries. The Terraform Cloud Providers implement most of the Clouds API, there is no need to reinvent the wheel.
* **[Pulumi](https://www.pulumi.com)**: This may be the major competitor of Terraform or Terranova. With Pulumi you can use many clouds using multiple programming languages. However, Terranova is a mix between Go and Terraform, if you like or use Terraform and want to use a programming language like Go, you have the same benefits as Pulumi and still have the benefits of using a DSL like HCL.
* **Terraform**: It's recommended to use Terraform instead of a Go program with Terranova when:
  * The code is simple
  * The code has to be easy to maintain
  * There is no need to add extra logic or HCL provide all the required functionalities
  * Use of external services such as Terraform Cloud or IBM Cloud Schematics

## Simple example

`main.go`

```go
package main

import (
	"log"
	"os"

	"github.com/johandry/terranova"
	"github.com/johandry/terranova/logger"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var code string

const stateFilename = "simple.tfstate"

func main() {
	count := 0
	keyName := "demo"

	log := log.New(os.Stderr, "", log.LstdFlags)
	logMiddleware := logger.NewMiddleware()
	defer logMiddleware.Close()

	platform, err := terranova.NewPlatform(code).
		SetMiddleware(logMiddleware).
		AddProvider("aws", aws.Provider()).
		Var("c", count).
		Var("key_name", keyName).
		PersistStateToFile(stateFilename)

	if err != nil {
		log.Fatalf("Fail to create the platform using state file %s. %s", stateFilename, err)
	}

	terminate := (count == 0)
	if err := platform.Apply(terminate); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}
}

func init() {
	code = `
  variable "c"    { default = 2 }
  variable "key_name" {}
  provider "aws" {
    region        = "us-west-2"
  }
  resource "aws_instance" "server" {
    instance_type = "t2.micro"
    ami           = "ami-6e1a0117"
    count         = "${var.c}"
    key_name      = "${var.key_name}"
  }
`
}
```



`go.mod`

```go
module github.com/johandry/terranova-examples/simple

go 1.15

require (
	github.com/johandry/terranova v0.0.4
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191003145700-f8707a46c6ec
)

replace github.com/terraform-providers/terraform-provider-tls => github.com/terraform-providers/terraform-provider-tls v1.2.1-0.20190816230231-0790c4b40281
```

Build and use

```bash
go build -o terractl .

./terractl
```

## Demo

The example in this directory creates a instance and deploy LetsChat, a simple web chat application.

After build or download the binary, execute:

```bash
# To view the Terraform code to execute:
letschat --status --ssh-access --export

# To create an instance with SSH Access
letschat --status --ssh-access --cloud ibm

# Destroy the instance
letschat --destroy
```
