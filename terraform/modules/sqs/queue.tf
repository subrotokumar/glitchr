resource "aws_sqs_queue" "storage_trigger_queue" {
  name =  var.storage_trigger_queue
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : "*",
        "Action" : "sqs:SendMessage",
        "Resource" : "${aws_sqs_queue.storage_trigger_queue.arn}",
        "Condition" : {
          "ArnEquals" : {
            "aws:SourceArn" : "${var.s3_bucket_arn}"
          }
        }
      }
    ]
  })
}
