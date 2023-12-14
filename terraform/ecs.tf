data "aws_ecr_image" "recipe_maker_image" {
  repository_name = aws_ecr_repository.recipe_maker_registry.name
  image_tag       = "latest"
}

data "aws_ecr_image" "recipe_maker_ui_image" {
  repository_name = aws_ecr_repository.recipe_maker_ui.name
  image_tag       = "latest"
}

resource "aws_ecs_cluster" "recipe_maker" {
  name = "recipe-maker"

  configuration {
    execute_command_configuration {
      logging = "OVERRIDE"

      log_configuration {
        cloud_watch_log_group_name = aws_cloudwatch_log_group.recipe_maker_log_group.name
      }
    }
  }

  depends_on = [aws_cloudwatch_log_group.recipe_maker_log_group]
}

resource "aws_ecs_cluster_capacity_providers" "recipe_maker" {
  cluster_name       = aws_ecs_cluster.recipe_maker.name
  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
  depends_on = [aws_ecs_cluster.recipe_maker]
}

resource "aws_ecs_task_definition" "recipe_maker_backend" {
  family                   = "recipe_maker"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  container_definitions = jsonencode([
    {
      name      = "api"
      image     = "${aws_ecr_repository.recipe_maker_registry.repository_url}:latest"
      cpu       = 128
      memory    = 256
      essential = true
      healthCheck = {
        command  = ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"]
        interval = 30
        timeout  = 5
        retries  = 3
      }
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.recipe_maker_log_group.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "recipe-maker"
        }
      }
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      environment = [
        {
          name  = "DEBUG_LOGGING"
          value = "true"
        },
        {
          name  = "VALIDATE_JWT"
          value = "false"
        },
        {
          name  = "JWT_SIGNING_SECRET"
          value = "somesigningsecret"
        },
        {
          name  = "DATABASE_URL"
          value = "postgres://${var.db_username}:${var.db_password}@${aws_db_instance.recipe_maker.endpoint}/${var.db_name}"
        },
      ]
    },
    {
      name      = "nginx"
      image     = "${aws_ecr_repository.recipe_maker_nginx.repository_url}:latest"
      cpu       = 128
      memory    = 256
      essential = true
      healthCheck = {
        command  = ["CMD-SHELL", "service nginx status"]
        interval = 30
        timeout  = 5
        retries  = 3
      }
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.recipe_maker_log_group.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "recipe-maker"
        }
      }
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
        },
        {
          containerPort = 443
          hostPort      = 443
        },
      ]
    },
  ])


  depends_on = [
    aws_ecs_cluster.recipe_maker,
    aws_iam_role.ecs_task_execution_role
  ]
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_ecs_task_definition" "recipe_maker_frontend" {
  family                   = "recipe_maker"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  container_definitions = jsonencode([
    {
      name      = "frontend"
      image     = "${aws_ecr_repository.recipe_maker_ui.repository_url}:latest"
      cpu       = 256
      memory    = 512
      essential = true
      healthCheck = {
        command  = ["CMD-SHELL", "curl -f http://localhost:3000/api/health || exit 1"]
        interval = 30
        timeout  = 5
        retries  = 3
      }
      environment = [
        {
          name  = "API_HOST"
          value = aws_lb.recipe_maker_lb.dns_name
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.recipe_maker_log_group.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "recipe-maker"
        }
      }
      portMappings = [
        {
          containerPort = 3000
          hostPort      = 3000
        },
      ]
    },
  ])

  depends_on = [
    aws_ecs_cluster.recipe_maker,
    aws_iam_role.ecs_task_execution_role
  ]
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecs-task-execution-role"
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
    "arn:aws:iam::aws:policy/AmazonRDSDataFullAccess",
    "arn:aws:iam::aws:policy/AmazonRDSFullAccess"
  ]

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ecs-tasks.amazonaws.com",
        },
      },
    ],
  })
}

resource "aws_ecs_service" "recipe_maker_api" {
  name                               = "recipe-maker-api"
  cluster                            = aws_ecs_cluster.recipe_maker.arn
  task_definition                    = aws_ecs_task_definition.recipe_maker_backend.arn
  desired_count                      = 2
  launch_type                        = "FARGATE"
  force_new_deployment               = true
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 50

  triggers = {
    redeployment = data.aws_ecr_image.recipe_maker_image.image_digest
  }

  network_configuration {
    subnets = [
      aws_subnet.public_subnet_1.id,
      aws_subnet.public_subnet_2.id,
    ]
    assign_public_ip = true
    security_groups  = [aws_security_group.alb_security_group.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.recipe_maker_target_group.arn
    container_name   = "nginx"
    container_port   = 80
  }

  depends_on = [
    aws_iam_role.ecs_task_execution_role,
    aws_ecs_cluster.recipe_maker,
    aws_lb_target_group.recipe_maker_target_group,
  ]
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_ecs_service" "recipe_maker_ui" {
  name                               = "recipe-maker-ui"
  cluster                            = aws_ecs_cluster.recipe_maker.arn
  task_definition                    = aws_ecs_task_definition.recipe_maker_frontend.arn
  desired_count                      = 2
  launch_type                        = "FARGATE"
  force_new_deployment               = true
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 50

  triggers = {
    redeployment = data.aws_ecr_image.recipe_maker_ui_image.image_digest
  }

  network_configuration {
    subnets = [
      aws_subnet.public_subnet_1.id,
      aws_subnet.public_subnet_2.id,
    ]
    assign_public_ip = true
    security_groups  = [aws_security_group.alb_security_group.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.recipe_maker_ui_tg.arn
    container_name   = "frontend"
    container_port   = 3000
  }

  depends_on = [
    aws_iam_role.ecs_task_execution_role,
    aws_ecs_cluster.recipe_maker,
    aws_lb_target_group.recipe_maker_ui_tg,
  ]
  tags = {
    service = "recipe-maker"
  }
}
