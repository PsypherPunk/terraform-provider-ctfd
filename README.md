# `terraform-provider-ctfd`

This repository is a built from the
[`terraform-provider-scaffolding`](https://github.com/hashicorp/terraform-provider-scaffolding)
GitHub [template](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-on-github/creating-a-repository-from-a-template).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/PsypherPunk/ctfd` to your Terraform provider:

```
go get github.com/PsypherPunk/ctfd
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

To say that these are early days ain't the half of it. I've never written a
Terraform provider before and Go and I don't really get on.

Nevertheless, so far:

- the provider assumes a running [CTFd](https://github.com/CTFd/CTFd) instance
  with the Admin. user already created: this is used to interact with the CTFd
  instance:

```hcl
provider "ctfd" {
  username = "admin"
  password = "admin"
  url      = "http://0.0.0.0:8000/"
}
```

- the only thing implemented are the *Challenges* as a
  [Data Source](https://www.terraform.io/docs/language/data-sources/index.html).

```hcl
data "ctfd_challenges" "challenges" {}

output "challenges" {
  value = data.ctfd_challenges.challenges
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (see
[Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
