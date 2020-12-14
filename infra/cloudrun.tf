resource "google_cloud_run_service" "storage_api" {
  name = "storage-api"
  location = "us-east4"

  template {
    spec {
      container_concurrency = 10
      containers {
        image = "gcr.io/winter-workshop/endpoints-runtime-serverless:2.21.0-storage-api-gh6esglzca-uk.a.run.app-2020-12-14r0"
      }
    }
  }

  depends_on = [
    google_project_service.run]
}

resource "google_cloud_run_service_iam_policy" "storage_gateway_noauth" {
  service = google_cloud_run_service.storage_api.name
  project = google_cloud_run_service.storage_api.project
  location = google_cloud_run_service.storage_api.location

  policy_data = data.google_iam_policy.noauth.policy_data
}

resource "google_cloud_run_service" "storage" {
  name = "storage-svc"
  location = "us-east4"

  template {
    spec {
      container_concurrency = 10
      containers {
        image = "us.gcr.io/winter-workshop/salzr/cloudrun-storage:v0.0.11"
        command = [
          "/app/server"]
        env {
          name = "STORAGE_PROJECT_ID"
          value = data.google_project.winter-workshop.project_id
        }
        env {
          name = "STORAGE_BUCKET"
          value = google_storage_bucket.cloudrun-storage.name
        }
      }
    }
  }

  traffic {
    percent = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.run]
}
