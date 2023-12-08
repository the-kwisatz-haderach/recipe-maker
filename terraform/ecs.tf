resource "aws_ecs_cluster" "recipe_maker" {
  name = "recipe-maker"
}

resource "aws_ecs_cluster_capacity_providers" "recipe_maker" {
  cluster_name = aws_ecs_cluster.recipe_maker.name

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
      image     = "${var.aws_account}.dkr.ecr.eu-north-1.amazonaws.com/recipe-maker:latest"
      cpu       = 256
      memory    = 512
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
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

resource "aws_lb" "recipe_maker_lb" {
  name               = "recipe-maker-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_security_group.id]
  subnets = [
    aws_subnet.main_subnet_1.id,
    aws_subnet.main_subnet_2.id,
    aws_subnet.main_subnet_3.id,
  ]

  enable_deletion_protection = false

  depends_on = [
    aws_security_group.alb_security_group,
    aws_subnet.main_subnet_1,
    aws_subnet.main_subnet_2,
    aws_subnet.main_subnet_3,
  ]
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_lb_target_group" "recipe_maker_lb_target_group" {
  name        = "recipe-maker-lb"
  target_type = "ip"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  depends_on  = [aws_vpc.main]
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.recipe_maker_lb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.recipe_maker_lb_target_group.arn
  }
  depends_on = [aws_lb_target_group.recipe_maker_lb_target_group]
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_ecs_service" "recipe_maker_api" {
  name            = "recipe-maker-api"
  cluster         = aws_ecs_cluster.recipe_maker.arn
  task_definition = aws_ecs_task_definition.recipe_maker.arn
  desired_count   = 1
  launch_type     = "FARGATE"
  network_configuration {
    subnets = [
      aws_subnet.main_subnet_1.id,
      aws_subnet.main_subnet_2.id,
      aws_subnet.main_subnet_3.id,
    ]
    security_groups = [aws_security_group.alb_security_group.id]
  }
  # iam_role        = aws_iam_role.foo.arn

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
