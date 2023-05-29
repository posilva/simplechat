resource "aws_s3_bucket" "deploy" {
  bucket = "${local.name}-deploy-${random_integer.suffix.result}"
  acl    = "private"
}

resource "random_integer" "suffix" {
  min = 100
  max = 999
}