data "digitalocean_project" "karavel" {
  name = "Karavel"
}

module "cluster" {
  source                 = "./cluster"
  cluster_name           = var.cluster_name
  cluster_version_prefix = "1.21."
  region                 = var.do_region
  project_id             = data.digitalocean_project.karavel.id
}

