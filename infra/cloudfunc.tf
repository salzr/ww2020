resource "google_cloudfunctions_function" "vision" {
  name = "cf-vision"
  description = "Cloud function in charge of image categorization"
  runtime = "go113"
  entry_point = "ProcessEvent"
  source_archive_bucket = google_storage_bucket.cloudfunc-vision-src.name
  source_archive_object = google_storage_bucket_object.cloudfunc-vision.name
  ingress_settings = "ALLOW_INTERNAL_ONLY"

  environment_variables = {
    PROJECT_ID = data.google_project.winter-workshop.name
  }

  event_trigger {
    resource = google_storage_bucket.cloudrun-storage.name
    event_type = "google.storage.object.finalize"
  }
}