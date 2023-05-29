locals {
  region      = var.region
  name        = var.service_name
  namespace   = var.namespace
  stage       = var.stage
  environment = var.environment

  vpc_cidr = "172.20.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  container_name = "${local.name}_demo"
  container_port = 8081

  tags = {
    Name       = local.name
    Owner      = "posilva@gmail.com"
    Repository = "https://github.com/posilva/${local.name}"
  }
}
