data "digitalocean_kubernetes_versions" "example" {
  version_prefix = var.cluster_version_prefix
}

resource "digitalocean_kubernetes_cluster" "cluster" {
  name   = var.cluster_name
  region = var.region
  version = data.digitalocean_kubernetes_versions.example.latest_version
  vpc_uuid = digitalocean_vpc.vpc.id

  maintenance_policy {
    start_time  = "04:00"
    day         = "sunday"
  }

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-4gb"
    node_count = 5
  }
}

resource "digitalocean_project_resources" "project_resources" {
  project = var.project_id
  resources = [
    digitalocean_kubernetes_cluster.cluster.urn
  ]
}
