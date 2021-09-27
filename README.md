# Fleet Terraform provider

This repo is a proof of concept of how a fleet provider for terraform could work

## Build provider

Run the following command to build the provider

```shell
$ make build
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```
