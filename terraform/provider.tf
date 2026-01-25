# Configure the Google provider
provider "google" {
  # The project ID where resources will be created.
  # This value can also be set via the GOOGLE_CLOUD_PROJECT environment variable.
  project = "api-project-119360632367"
  # The region to manage resources in (optional, but recommended).
  region = "us-east1"
}