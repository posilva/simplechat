module "security_group" {
  source  = "cloudposse/security-group/aws"
  version = "2.1.0"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  allow_all_egress = true

  rules = [
    {
      key                      = "HTTP"
      type                     = "ingress"
      from_port                = 80
      to_port                  = 80
      protocol                 = "tcp"
      source_security_group_id = aws_security_group.lb_sg.id
      description              = "Allow HTTP from load balancer"
    }
  ]
  vpc_id = module.vpc.vpc_id
}
