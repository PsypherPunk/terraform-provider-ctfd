terraform {
  required_providers {
    ctfd = {
      source  = "psypherpunk.io/ctfd/ctfd"
      version = "0.1.0"
    }
  }
}

provider "ctfd" {
  username = "admin"
  password = "admin"
  url      = "http://0.0.0.0:8000/"
}
