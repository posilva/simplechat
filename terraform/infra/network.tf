
data "aws_availability_zones" "available" {}

module "vpc" {
  source  = "cloudposse/vpc/aws"
  version = "2.1.0"

  enabled = true

  name        = local.name
  namespace   = local.namespace
  stage       = local.stage
  environment = local.environment

  ipv4_primary_cidr_block = local.vpc_cidr

  internet_gateway_enabled = true

  dns_hostnames_enabled = true
  dns_support_enabled   = true
}

module "subnets" {
  source  = "cloudposse/dynamic-subnets/aws"
  version = "2.0.4"

  namespace   = local.namespace
  stage       = local.stage
  environment = local.environment

  vpc_id = module.vpc.vpc_id

  availability_zones = local.azs

  igw_id          = [module.vpc.igw_id]
  ipv4_cidr_block = [module.vpc.vpc_cidr_block]

  nat_gateway_enabled  = true
  nat_instance_enabled = true

  private_subnets_enabled = false
  public_subnets_enabled  = true

}
