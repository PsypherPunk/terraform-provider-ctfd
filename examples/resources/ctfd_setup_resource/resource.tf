resource "ctfd_setup" "setup" {
  name               = "Test CTFd!"
  description        = "Example CTFd setup."
  admin_email        = "admin@example.com"
  configuration_path = "/tmp/juice-shop-ctf.zip"

  email {
    username     = "admin@example.com"
    password     = "secure-password"
    from_address = "admin+ctfd@example.com"
    server       = "smtp.example.com"
    port         = 587
    use_auth     = true
    use_tls      = true
    use_ssl      = false
  }
}

