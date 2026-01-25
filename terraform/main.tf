locals {
  project    = "api-project-119360632367"
  region     = "us-east1"
  repository = "calendar"
  image_name = "main"
  tag        = "qa"
}
resource "google_cloud_run_v2_service" "default" {
  name                 = "calendar"
  location             = local.region
  deletion_protection  = false
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true

  scaling {
    max_instance_count = 10
  }

  template {
    containers {
      image = data.google_artifact_registry_docker_image.my_image.self_link
    }
  }
}

data "google_artifact_registry_docker_image" "my_image" {
  location      = local.region
  repository_id = local.repository
  image_name    = var.image
}