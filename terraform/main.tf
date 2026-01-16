resource "google_cloud_run_v2_service" "default" {
  name     = "calendar"
  location = "us-east1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"

  scaling {
    max_instance_count = 10
  }

  template {
    containers {
      image = "us-east1-docker.pkg.dev/api-project-119360632367/family-cal/calendar"
    }
  }
}