module "config" {
  source="./config"
}

provider "google" {
  credentials = "${file("./config/google_cloud_credentials.json")}"
  project = "${module.config.google_project_id}"
  region = "${module.config.region}"
}

resource "google_compute_global_address" "static" {
  name = "ipv4-address"
}

resource "google_dns_managed_zone" "myzone" {
  name        = "myzone"
  dns_name    = "${module.config.domain}."

  dnssec_config {
   state      = "on"
  }
}

resource "google_dns_record_set" "pullassistant_com_set" {
  name = "${google_dns_managed_zone.myzone.dns_name}"
  managed_zone = "${google_dns_managed_zone.myzone.name}"
  type = "A"
  ttl  = 300

  rrdatas = ["185.199.108.153", "185.199.109.153", "185.199.110.153", "185.199.111.153"]
}

resource "google_dns_record_set" "pullassistant_com_mx" {
  name = "${google_dns_managed_zone.myzone.dns_name}"
  managed_zone = "${google_dns_managed_zone.myzone.name}"
  type = "MX"
  ttl  = 300

  rrdatas = ["10 mx1.improvmx.com.", "20 mx2.improvmx.com."]
}

resource "google_dns_record_set" "app_pullassistant_com_set" {
  name = "app.${google_dns_managed_zone.myzone.dns_name}"
  managed_zone = "${google_dns_managed_zone.myzone.name}"
  type = "A"
  ttl  = 300

  rrdatas = ["${google_compute_global_address.static.address}"]
}

resource "google_container_cluster" "gke_cluster" {
  name               = "gke-cluster"
  location           = "${module.config.region}"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count = 1

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

resource "google_container_node_pool" "gke_nodes" {
  name       = "gke-nodes"
  location   = "${module.config.region}"
  cluster    = "${google_container_cluster.gke_cluster.name}"
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "n1-standard-1"

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/devstorage.read_only"
    ]
  }
}
