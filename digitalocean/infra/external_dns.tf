resource "aws_iam_role" "external_dns" {
  name = "${var.cluster_name}-external-dns"

  inline_policy {
    name = "Route53ReadWriteAccess"
    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          "Action" = "route53:ChangeResourceRecordSets",
          "Resource" = "arn:aws:route53:::hostedzone/${var.karavel_io_hosted_zone}",
          "Effect" = "Allow"
        },
        {
          "Action" = [
            "route53:ListHostedZones",
            "route53:ListResourceRecordSets"
          ],
          "Resource" = "*",
          "Effect" = "Allow"
        }]
    })
  }
  assume_role_policy = jsonencode({
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          AWS = aws_iam_user.kube2iam.arn
        }
      },
    ]
  })
}

output "external_dns_iam_role" {
  value = aws_iam_role.external_dns.arn
}
