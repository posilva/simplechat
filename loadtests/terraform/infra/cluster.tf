resource "aws_iam_policy" "allow_s3_code_deploy" {
  name        = "${local.name}-allow-s3-code-deploy"
  path        = "/"
  description = "Allow to get code deploy from s3 bucket"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:GetObject",
          "s3:ListBucket"
        ]
        Effect = "Allow"
        Resource = [
          aws_s3_bucket.deploy.arn,
          "${aws_s3_bucket.deploy.arn}"
        ]
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "s3_codedeploy" {
  role       = aws_iam_role.screamer.name
  policy_arn = aws_iam_policy.allow_s3_code_deploy.arn
}


resource "aws_iam_role_policy_attachment" "instance_profile_codedeploy" {
  # access to S3 to CodeDeploy managed instances
  role       = aws_iam_role.screamer.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy"
}

resource "aws_iam_role_policy_attachment" "ssm-policy" {
  role       = aws_iam_role.screamer.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "screamer" {
  name = "${local.name}-instance_profile"
  role = aws_iam_role.screamer.name
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"
      identifiers = [
        "ec2.amazonaws.com"
      ]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "screamer" {
  name               = "${local.name}-instance-role"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}


resource "aws_security_group_rule" "screamer_self" {
  protocol          = -1
  from_port         = 0
  to_port           = 65535
  type              = "ingress"
  self              = true
  security_group_id = aws_security_group.screamer_cluster.id
}
resource "aws_security_group_rule" "outbound" {
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  type              = "egress"
  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.screamer_cluster.id
}

resource "aws_security_group" "screamer_cluster" {
  name        = "${local.name}-intra-cluster"
  description = "Allow TLS inbound traffic"
  vpc_id      = module.vpc.vpc_id
}

data "aws_ami" "screamer" {

  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "architecture"
    values = ["arm64"]

  }
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm*"]
  }
}

resource "aws_launch_template" "screamer" {
  name_prefix   = "screamer"
  image_id      = data.aws_ami.screamer.image_id
  instance_type = "t3g.micro"

  iam_instance_profile {
    name = aws_iam_instance_profile.screamer.name
  }
  monitoring {
    enabled = true
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = ["${aws_security_group.screamer_cluster.id}"]
  }

  user_data              = base64encode(templatefile("codedeploy_agent_install.sh", { aws_region = local.region }))
  update_default_version = true
}

resource "aws_autoscaling_group" "screamer" {
  name                = "${local.name}-asg"
  capacity_rebalance  = true
  desired_capacity    = 1
  max_size            = 2
  min_size            = 1
  vpc_zone_identifier = module.subnets.public_subnet_ids

  mixed_instances_policy {
    instances_distribution {
      on_demand_base_capacity                  = 0
      on_demand_percentage_above_base_capacity = 0
      spot_allocation_strategy                 = "capacity-optimized"
    }

    launch_template {
      launch_template_specification {
        launch_template_id = aws_launch_template.screamer.id
        version            = aws_launch_template.screamer.default_version
      }

      override {
        instance_type     = "t4g.medium"
        weighted_capacity = "3"
      }

      override {
        instance_type     = "t4g.small"
        weighted_capacity = "2"
      }
    }
  }
  instance_refresh {
    strategy = "Rolling"
    preferences {
      min_healthy_percentage = 50
    }
    triggers = ["tag"]
  }

  tag {
    key                 = "version"
    value               = "2.2"
    propagate_at_launch = true
  }
  tag {
    key                 = "Name"
    value               = "${local.name}-asg"
    propagate_at_launch = true
  }

}
