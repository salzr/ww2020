resource "google_project_service" "firebase" {
  project = data.google_project.winter-workshop.name
  service = "firebase.googleapis.com"
}

resource "google_project_service" "gcr" {
  project = data.google_project.winter-workshop.name
  service = "containerregistry.googleapis.com"
}

resource "google_project_service" "run" {
  project = data.google_project.winter-workshop.name
  service = "run.googleapis.com"
}

resource "google_project_service" "cloudbuild" {
  project = data.google_project.winter-workshop.name
  service = "cloudbuild.googleapis.com"
}

resource "google_project_service" "servicecontrol" {
  project = data.google_project.winter-workshop.name
  service = "servicecontrol.googleapis.com"
}

resource "google_project_service" "endpoints" {
  project = data.google_project.winter-workshop.name
  service = "endpoints.googleapis.com"
}

resource "google_project_service" "storage-api" {
  project = data.google_project.winter-workshop.name
  service = google_endpoints_service.storage.service_name
}

resource "google_project_service" "vision" {
  project = data.google_project.winter-workshop.name
  service = "vision.googleapis.com"
}