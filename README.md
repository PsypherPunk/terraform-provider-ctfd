# `terraform-provider-ctfd`

This repository is a built from the
[`terraform-provider-scaffolding`](https://github.com/hashicorp/terraform-provider-scaffolding)
GitHub [template](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-on-github/creating-a-repository-from-a-template).

The contents of the `api` package are taken from the
[`hashicups-client-go`](https://github.com/hashicorp-demoapp/hashicups-client-go)
demo. API client.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```sh
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using
Go modules.

To add a new dependency `github.com/PsypherPunk/ctfd` to your Terraform provider:

```sh
go get github.com/PsypherPunk/ctfd
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

The provider assumes a running [CTFd](https://github.com/CTFd/CTFd) instance
with the Admin. user already created: this is used to interact with the CTFd
instance:

```hcl
provider "ctfd" {
  username = "admin"
  password = "admin"
  url      = "http://0.0.0.0:8000/"
}
```

### [Data Sources](https://www.terraform.io/docs/language/data-sources/index.html)

#### Challenges

```hcl
data "ctfd_challenges" "challenges" {}

output "challenges" {
  value = data.ctfd_challenges.challenges
}
```

### [Resources](https://www.terraform.io/docs/language/resources/index.html)

#### Setup

i.e. the initial configuration of the Admin. user and import of challenges from
a [`juice-shop-ctf`](https://github.com/juice-shop/juice-shop-ctf)) instance:

```hcl
resource "ctfd_setup" "setup" {
  name               = "CTFd Instance"
  description        = "This is a test CTFd intance."
  admin_email        = "ctfd.admin@example.com"
  configuration_path = "/tmp/juice-shop-ctf-export.zip"
}
```

#### Teams

```hcl
resource "ctfd_team" "first_team" {
  name     = "First Team"
  email    = "first.team@example.com"
  password = "pass"
}
```

#### Users

```hcl
resource "ctfd_user" "first_user" {
  name     = "First User"
  email    = "first.user@example.com"
  password = "pass"
  type     = "user"
}
```

#### Membership

```hcl
resource "ctfd_user_team_membership" "first_user" {
  team_id = ctfd_team.first_team.id
  user_id = ctfd_user.first_user.id
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
make testacc
```
