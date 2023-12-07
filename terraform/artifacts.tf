resource "aws_ecr_repository" "recipe_maker_registry" {
  name                 = "recipe-maker"
  image_tag_mutability = "MUTABLE"
  force_delete         = true
  tags = {
    tag-key = "recipe-maker"
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
          Federated = "arn:aws:iam::044984945511:oidc-provider/token.actions.githubusercontent.com"
        }
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:the-kwisatz-haderach/recipe-maker:*",
          },
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com",
          }
        }
      },
    ]
  })
  tags = {
    tag-key = "recipe-maker"
  }
}

data "aws_iam_policy_document" "ecr_policy_document" {
  version = "2012-10-17"

  statement {
    effect = "Allow"
    actions = [
      "ecr:CompleteLayerUpload",
      "ecr:GetAuthorizationToken",
      "ecr:UploadLayerPart",
      "ecr:InitiateLayerUpload",
      "ecr:BatchCheckLayerAvailability",
      "ecr:PutImage",
    ]
    resources = [aws_ecr_repository.recipe_maker_registry.arn]
  }
  depends_on = [aws_ecr_repository.recipe_maker_registry]
}

resource "aws_iam_policy" "ecr_access_policy" {
  name   = "ECRAccessPolicy"
  policy = data.aws_iam_policy_document.ecr_policy_document.json
  tags = {
    tag-key = "recipe-maker"
  }
}

resource "aws_iam_role_policy_attachment" "github_actions_ecr_policy_attachment" {
  role       = aws_iam_role.github_actions_role.name
  policy_arn = aws_iam_policy.ecr_access_policy.arn
}
