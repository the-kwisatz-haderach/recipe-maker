variable "region" {
  description = "default region used"
  type        = string
  default     = "eu-north-1"
}

variable "aws_account" {
  description = "aws account"
  type        = string
  default     = "044984945511"
}

variable "github_account" {
  description = "github account/org"
  type        = string
  default     = "the-kwisatz-haderach"
}

variable "github_repo" {
  description = "github repository"
  type        = string
  default     = "recipe-maker"
}
