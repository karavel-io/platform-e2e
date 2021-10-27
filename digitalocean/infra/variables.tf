variable "keybase_user" {
  type = string
}

variable "aws_region" {
  type = string
  default = "eu-west-1"
}

variable "do_region" {
  type = string
  default = "fra1"
}

variable "cluster_name" {
  type = string
  default = "karavel-e2e-testing"
}

variable "karavel_io_hosted_zone" {
  type = string
  default = "Z02598523M9WR611ST687"
}
