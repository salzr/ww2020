resource "google_project_service" "firebase" {
  project = data.google_project.winter-workshop.name
  service = "firebase.googleapis.com"
}

resource "google_project_service" "gcr" {
  project = data.google_project.winter-workshop.name
  service = "containerregistry.googleapis.com"
}