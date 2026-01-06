module "cognito" {
  source = "../modules/cognito"

  user_pool   = var.user_pool
  client_name = var.user_pool_client
}

module "processed_media_s3" {
  source = "../modules/s3"

  bucket = var.processed_media_bucket
}


module "raw_media_s3" {
  source = "../modules/s3"

  bucket = var.raw_media_bucket
}

module "sqs" {
  source = "../modules/sqs"

  storage_trigger_queue = var.storage_trigger_queue
  s3_bucket_arn = module.raw_media_s3.bucket_arn
}

module "s3_notification" {
  source = "../modules/s3_notification"

  queue_arn = module.sqs.storage_trigger_queue_arn
  s3_bucket_id = module.raw_media_s3.bucket_arn
}