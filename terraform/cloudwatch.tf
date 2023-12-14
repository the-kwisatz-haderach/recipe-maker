resource "aws_cloudwatch_log_group" "recipe_maker_log_group" {
  name              = "recipe-maker-log-group"
  retention_in_days = 7
}
