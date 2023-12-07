resource "aws_s3_bucket" "recipe_maker_bucket" {
  bucket        = "recipe-maker"
  force_destroy = true
  tags = {
    tag-key = "recipe-maker"
  }
}

resource "aws_s3_bucket_ownership_controls" "recipe_maker" {
  bucket = aws_s3_bucket.recipe_maker_bucket.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "recipe_maker" {
  depends_on = [aws_s3_bucket_ownership_controls.recipe_maker]

  bucket = aws_s3_bucket.recipe_maker_bucket.id
  acl    = "private"
}
