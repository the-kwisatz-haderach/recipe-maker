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
