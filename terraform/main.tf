# UGC Ownership Table
module "dynamodb_table_simple_chat" {
  source  = "cloudposse/dynamodb/aws"
  version = "0.32.0"

  namespace    = var.namespace
  stage        = var.stage
  environment  = var.environment
  name         = "simple_chat"
  hash_key     = "pk"
  range_key    = "sk"
  billing_mode = var.dynamodb_billing_mode

  global_secondary_index_map = [
    {
      name            = "gsi-source-group-idx"
      hash_key        = "source"
      range_key       = "pk"
      projection_type = "ALL"
      # capacities are not required for GSI but the schema requires them
      read_capacity  = null
      write_capacity = null
      # non_key_attributes must be set to null if the projection_type is "ALL"
      non_key_attributes = null
    }
  ]

  dynamodb_attributes = [
    {
      name = "pk"
      type = "S"
    },
    {
      name = "sk"
      type = "S"
    },
    {
      name = "source"
      type = "S"
    }
  ]
}
