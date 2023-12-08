resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_subnet" "main_subnet_1" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "${var.region}a"
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_subnet" "main_subnet_2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "${var.region}b"
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_subnet" "main_subnet_3" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "${var.region}c"
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_security_group" "alb_security_group" {
  vpc_id      = aws_vpc.main.id
  description = "Security group for Application Load Balancer"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1" # Allow all outbound traffic
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_internet_gateway" "main_igw" {
  vpc_id = aws_vpc.main.id
}


resource "aws_route_table" "main_route_table" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route" "route" {
  route_table_id         = aws_route_table.main_route_table.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main_igw.id
}

resource "aws_route_table_association" "route_table_association_a" {
  subnet_id      = aws_subnet.main_subnet_1.id
  route_table_id = aws_route_table.main_route_table.id
}

resource "aws_route_table_association" "route_table_association_b" {
  subnet_id      = aws_subnet.main_subnet_2.id
  route_table_id = aws_route_table.main_route_table.id
}

resource "aws_route_table_association" "route_table_association_c" {
  subnet_id      = aws_subnet.main_subnet_3.id
  route_table_id = aws_route_table.main_route_table.id
}
