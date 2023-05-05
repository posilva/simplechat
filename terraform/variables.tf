# Global
#
variable "region" {
  type        = string
  description = "AWS region"
  default     = "us-east-1"
}

variable "namespace" {
  type        = string
  description = "Namespace (e.g. `local`)"
  default     = "local"
}

variable "stage" {
  type        = string
  description = "Stage (e.g. `prod`, `dev`, `staging`)"
  default     = "dev"
}

variable "environment" {
  type        = string
  description = "Environment (e.g. `prod`, `dev`, `staging`)"
  default     = "dev"
}

# service's name
variable "service_name" {
  type        = string
  description = "Service name"
  default     = "simple_chat"
}

# services's port (as exposed in the container)
variable "service_port" {
  type        = number
  description = "Service port"
  default     = 18808
}

variable "aws_account_id" {
  type        = string
  description = "AWS account ID"
  default     = "123456789012"
}

# DynamoDB
#
variable "dynamodb_billing_mode" {
  type        = string
  description = "DynamoDB billing mode"
  default     = "PAY_PER_REQUEST"
}
