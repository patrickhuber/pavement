// requires GOOGLE_CLOUD_KEYFILE_JSON environment variable set
provider "google" {
  project = "fe-phuber"
  zone    = "us-east1-b"
  region  = "us-east1"
}

resource "google_compute_project_metadata_item" "oslogin" {
  key   = "enable-oslogin"
  value = "TRUE"
}

variable "username" {
  default = "pavement"
}
