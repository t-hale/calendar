# Configure the Google provider
provider "google" {
  # The project ID where resources will be created.
  # This value can also be set via the GOOGLE_CLOUD_PROJECT environment variable.
  project = "api-project-119360632367"
  # The region to manage resources in (optional, but recommended).
  # region = "us-central1"
}

resource "google_cloud_run_v2_service" "default" {
  name     = "calendar"
  location = "us-east1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true

  scaling {
    max_instance_count = 10
  }

  template {
    containers {
      image = "us-east1-docker.pkg.dev/api-project-119360632367/calendar/main"
    }
  }
}