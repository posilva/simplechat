module "ecr" {
  source  = "cloudposse/ecr/aws"
  version = "0.37.0"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  use_fullname = true

  image_tag_mutability    = "MUTABLE"
  enable_lifecycle_policy = true

}
