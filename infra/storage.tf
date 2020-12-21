resource "google_storage_bucket" "cloudrun-storage" {
  name = format("cloudrun-storage-%s", data.google_project.winter-workshop.number)
  force_destroy = true
}

resource "google_storage_bucket" "cloudfunc-vision-src" {
  name = "cloudfunc-vision-src"
}

resource "google_storage_bucket_object" "cloudfunc-vision" {
  name = "vision.zip"
  bucket = google_storage_bucket.cloudfunc-vision-src.name
  source = pathexpand("../cloudfunc/vision/vision.zip")
}