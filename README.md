Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-packet`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-packet
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-packet
$ make build
```

Using the provider
----------------------

The packet provider will be installed on `terraform init` of a template using any of the `packet_*` resources.

Available resource and datasources are documented at [https://www.terraform.io/docs/providers/packet/index.html](https://www.terraform.io/docs/providers/packet/index.html).


Developing the provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-packet
...
```


Testing provider code
---------------------------

We have mostly acceptance tests in the provider. There's no point for you to run them all, but you should run the one covering the functionality which you change. The acceptance test run will cost you some money, so feel free to abstain. The acceptance test suite will be run for your PR during the review process.

To run an acceptance test, find the relevant test function in `*_test.go` (for example TestAccPacketDevice_Basic), and run it as

```sh
TF_ACC=1 go test -v -timeout=20m -run=TestAccPacketDevice_Basic
```

If you want to see HTTP traffic, set `TF_LOG=DEBUG`, i.e.

```sh
TF_LOG=DEBUG TF_ACC=1 go test -v -timeout=20m -run=TestAccPacketDevice_Basic
```



Testing the provider with Terraform
---------------------------------------

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider) above), it can be incorporated within your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will then use the plugin from your development environment.

```sh
$ terraform init -plugin-dir $GOPATH/bin
```
