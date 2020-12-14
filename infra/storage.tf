resource "google_storage_bucket" "cloudrun-storage" {
  name = format("cloudrun-storage-%s", data.google_project.winter-workshop.number)
  force_destroy = true
}