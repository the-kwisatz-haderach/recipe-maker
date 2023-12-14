resource "aws_db_subnet_group" "recipe_maker" {
  name = "recipe-maker"
  subnet_ids = [
    aws_subnet.public_subnet_1.id,
    aws_subnet.public_subnet_2.id,
  ]

  tags = {
    service = "recipe-maker"
  }
}

resource "aws_db_parameter_group" "recipe_maker" {
  name   = "recipe-maker"
  family = "postgres15"

  parameter {
    name  = "log_connections"
    value = "1"
  }
}

resource "aws_db_instance" "recipe_maker" {
  identifier             = "recipe-maker"
  instance_class         = "db.t3.micro"
  allocated_storage      = 5
  engine                 = "postgres"
  db_name                = var.db_name
  engine_version         = "15.4"
  port                   = var.db_port
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.recipe_maker.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.recipe_maker.name
  publicly_accessible    = false
  skip_final_snapshot    = true

  depends_on = [
    aws_db_parameter_group.recipe_maker,
    aws_db_subnet_group.recipe_maker
  ]
}

