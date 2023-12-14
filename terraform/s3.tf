resource "aws_s3_bucket" "recipe_maker_bucket" {
  bucket        = "recipe-maker"
  force_destroy = true
  tags = {
    service = "recipe-maker"
  }
  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_ownership_controls" "recipe_maker" {
  bucket = aws_s3_bucket.recipe_maker_bucket.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_acl" "recipe_maker" {
  depends_on = [aws_s3_bucket_ownership_controls.recipe_maker]

  bucket = aws_s3_bucket.recipe_maker_bucket.id
  acl    = "private"
  lifecycle {
    prevent_destroy = true
  }
}
