resource "aws_ecr_repository" "recipe_maker_registry" {
  name                 = "recipe-maker"
  image_tag_mutability = "MUTABLE"
  force_delete         = true
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_iam_role" "github_actions_role" {
  name = "github_actions_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRoleWithWebIdentity"
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::${var.aws_account}:oidc-provider/token.actions.githubusercontent.com"
        }
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${var.github_account}/${var.github_repo}:*",
          },
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com",
          }
        }
      },
    ]
  })
  tags = {
    service = "recipe-maker"
  }
}

data "aws_iam_policy_document" "ecr_policy_document" {
  version = "2012-10-17"

  statement {
    effect = "Allow"
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:BatchGetImage",
      "ecr:CompleteLayerUpload",
      "ecr:GetDownloadUrlForLayer",
      "ecr:InitiateLayerUpload",
      "ecr:PutImage",
      "ecr:UploadLayerPart"
    ]
    resources = [aws_ecr_repository.recipe_maker_registry.arn]
  }
  statement {
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }

  depends_on = [aws_ecr_repository.recipe_maker_registry]
}

resource "aws_iam_policy" "ecr_access_policy" {
  name   = "ECRAccessPolicy"
  policy = data.aws_iam_policy_document.ecr_policy_document.json
  tags = {
    service = "recipe-maker"
  }
}

resource "aws_iam_role_policy_attachment" "github_actions_ecr_policy_attachment" {
  role       = aws_iam_role.github_actions_role.name
  policy_arn = aws_iam_policy.ecr_access_policy.arn
}

resource "aws_ecr_lifecycle_policy" "remove_untagged" {
  repository = aws_ecr_repository.recipe_maker_registry.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1,
        description  = "Expire images older than 7 days",
        selection = {
          tagStatus   = "untagged",
          countType   = "sinceImagePushed",
          countUnit   = "days",
          countNumber = 7
        },
        action = {
          type = "expire"
        }
      }
    ]
  })
}
