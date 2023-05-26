locals {
  region      = var.region
  name        = var.service_name
  namespace   = var.namespace
  stage       = var.stage
  environment = var.environment

  container_name = "${local.name}_demo"
  container_port = 8081

  tags = {
    Name       = local.name
    Owner      = "posilva@gmail.com"
    Repository = "https://github.com/posilva/${local.name}"
  }
}

data "aws_ecr_repository" "service" {
  name = "${local.namespace}-${local.environment}-${local.stage}-${local.name}"

}
resource "null_resource" "docker_packaging" {

  provisioner "local-exec" {
    command = <<EOF
    aws ecr get-login-password --region ${local.region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com
    cd ..
    cd ..
    docker buildx build --no-cache --load --platform linux/arm64 -t "${data.aws_ecr_repository.service.repository_url}:latest" -f Dockerfile .
    docker push "${data.aws_ecr_repository.service.repository_url}:latest"
    EOF
  }

  triggers = {
    "run_at" = timestamp()
  }
}