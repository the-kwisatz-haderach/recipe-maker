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

variable "db_password" {
  description = "root user password for db"
  type        = string
  sensitive   = true
}

variable "db_port" {
  description = "port used by db"
  type        = string
  sensitive   = true
}

variable "db_username" {
  description = "root username for db"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "root username for db"
  type        = string
  default     = "recipe_maker"
}

variable "pgadmin_default_email" {
  description = "default email for pgadmin"
  type        = string
  sensitive   = true
}

variable "pgadmin_default_password" {
  description = "default password for pgadmin"
  type        = string
  sensitive   = true
}

variable "github_repositories" {
  description = "github repositories allowed identity federation"
  type        = set(string)
  default     = ["recipe-maker-ui", "recipe-maker"]
}
