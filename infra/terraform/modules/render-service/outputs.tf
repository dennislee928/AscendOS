output "service_id" {
  description = "Render service ID"
  value       = render_web_service.this.id
}

output "service_url" {
  description = "Public URL of the deployed Render service"
  value       = render_web_service.this.service_details.url
}

output "service_name" {
  description = "Name of the Render service"
  value       = render_web_service.this.name
}
