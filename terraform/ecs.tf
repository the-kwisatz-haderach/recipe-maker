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

resource "aws_ecs_task_definition" "recipe_maker" {
  family                   = "recipe_maker_api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  container_definitions = jsonencode([
    {
      name      = "recipe-maker-api"
      image     = "${aws_ecr_repository.recipe_maker_registry.repository_url}:latest"
      cpu       = 256
      memory    = 512
      essential = true
      healthCheck = {
        command  = ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"]
        interval = 30
        timeout  = 5
        retries  = 3
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
  name            = "recipe-maker-api"
  cluster         = aws_ecs_cluster.recipe_maker.arn
  task_definition = aws_ecs_task_definition.recipe_maker.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets = [
      aws_subnet.public_subnet_1.id,
      aws_subnet.public_subnet_2.id,
      aws_subnet.public_subnet_3.id,
    ]
    assign_public_ip = true
    security_groups  = [aws_security_group.alb_security_group.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.recipe_maker_lb_target_group.arn
    container_name   = "recipe-maker-api"
    container_port   = 8080
  }

  depends_on = [
    aws_iam_role.ecs_task_execution_role,
    aws_ecs_cluster.recipe_maker,
    aws_lb_target_group.recipe_maker_lb_target_group,
    aws_lb_listener.front_end
  ]
  tags = {
    service = "recipe-maker"
  }
}
