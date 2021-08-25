terraform {
  required_providers {
    ctfd = {
      source = "psypherpunk.io/ctfd/ctfd"
    }
  }
}

provider "ctfd" {
  username = "admin"
  password = "admin"
  url      = "http://0.0.0.0:8000/"
}
