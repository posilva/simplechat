provider "aws" {
  region = local.region
}
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.0.0"
    }
  }
}
// Used by get the current aws number account.
data "aws_caller_identity" "current" {
}