data "google_project" "winter-workshop" {
  project_id = "winter-workshop"
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers"]
  }
}