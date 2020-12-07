resource "google_firebase_project" "auth" {
  provider = google-beta
  project = data.google_project.winter-workshop.number
  depends_on = [
    google_project_service.firebase]
}

resource "google_firebase_web_app" "ww2020" {
  provider = google-beta
  project = data.google_project.winter-workshop.number
  display_name = "Winter Workshop 2020"
  depends_on = [
    google_firebase_project.auth]
}