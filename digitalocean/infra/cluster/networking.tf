resource "digitalocean_vpc" "vpc" {
  name   = "${var.cluster_name}-network"
  region = var.region
}
