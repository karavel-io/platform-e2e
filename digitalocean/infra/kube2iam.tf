resource "aws_iam_user" "kube2iam" {
  name = "karavel-e2e-do-kube2iam"
}

resource "aws_iam_access_key" "kube2iam" {
  user    = aws_iam_user.kube2iam.name
  pgp_key = "keybase:${var.keybase_user}"
}

resource "aws_iam_user_policy_attachment" "attach_admin_access" {
  user       = aws_iam_user.kube2iam.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

output "kube2iam_access_key" {
  value = aws_iam_access_key.kube2iam.id
}

output "kube2iam_secret_key" {
  value = aws_iam_access_key.kube2iam.encrypted_secret
}
