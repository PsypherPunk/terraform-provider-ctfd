resource "ctfd_setup" "setup" {
  name               = "Test CTFd!"
  description        = "Example CTFd setup."
  admin_email        = "admin@example.com"
  configuration_path = "/tmp/juice-shop-ctf.zip"
}

