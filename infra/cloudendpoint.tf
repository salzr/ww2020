resource "google_endpoints_service" "storage" {
  service_name = trimprefix(google_cloud_run_service.storage_api.status[0]["url"], "https://")
  project = data.google_project.winter-workshop.name
  openapi_config = file("storage_endpoints.yaml")
}