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


  deploy_s3_bucket = "${local.name}-deploy"
  object_source    = "${path.module}/appdeploy.zip"
}

data "archive_file" "source" {
  type        = "zip"
  source_dir  = "./appdeploy"
  output_path = local.object_source

}

data "aws_s3_bucket" "selected" {
  bucket = local.deploy_s3_bucket

}

resource "null_resource" "upload_to_s3" {
  triggers = {
    md5 = data.archive_file.source.output_md5
  }


  provisioner "local-exec" {
    command = "aws s3 cp ${local.object_source} s3://${local.deploy_s3_bucket}/${local.name}-deploy-${data.archive_file.source.output_md5}/appdeploy.zip"
  }
}

