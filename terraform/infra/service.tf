
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

### Task execution role setup
data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ddb_policy" {
  statement {
    actions = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Query"
    ]
    resources = [
      "*" # TODO: We may lock this later 
    ]
  }
}

resource "aws_iam_policy" "ddb_policy" {
  name   = "${local.name}-ddb-policy"
  policy = data.aws_iam_policy_document.ddb_policy.json
}

resource "aws_iam_role_policy_attachment" "ecs-task-ddb-policy-attachment" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ddb_policy.arn
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name               = "${local.name}-execution-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
  tags               = local.tags
}

resource "aws_iam_role" "ecs_task_role" {
  name               = "${local.name}-task-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
  tags               = local.tags
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_ecs_task_definition" "main" {
  family = "${local.name}-family"

  network_mode = "awsvpc"
  runtime_platform {
    cpu_architecture        = "ARM64"
    operating_system_family = "LINUX"
  }
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  container_definitions = jsonencode([{
    name      = "${local.name}-container-${local.environment}"
    image     = "nginxdemos/hello"
    essential = true
    portMappings = [{
      protocol      = "tcp"
      containerPort = 80
      hostPort      = 80
    }]
  }])
}

resource "aws_ecs_service" "main" {
  name                               = "${local.name}-service-${local.environment}"
  cluster                            = aws_ecs_cluster.main.id
  task_definition                    = aws_ecs_task_definition.main.arn
  desired_count                      = 2
  deployment_minimum_healthy_percent = 50
  deployment_maximum_percent         = 200
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"


  network_configuration {
    security_groups  = [module.security_group.id]
    subnets          = module.subnets.public_subnet_ids
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.service.arn
    container_name   = "${local.name}-container-${local.environment}"
    container_port   = 80
  }

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }
}

resource "aws_lb_target_group" "service" {
  name        = "${local.name}-tg"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"
}

resource "aws_lb_listener_rule" "service_rule" {
  listener_arn = aws_lb_listener.lb_listener.arn
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.service.arn
  }
  condition {
    path_pattern {
      values = ["/*"]
    }
  }
}
