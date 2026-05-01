output "supabase_project_id" {
  description = "Supabase project ID"
  value       = supabase_project.main.id
}

output "supabase_api_url" {
  description = "Supabase REST API URL"
  value       = supabase_project.main.api_url
}

output "upstash_redis_endpoint" {
  description = "Upstash Redis hostname:port"
  value       = "${upstash_redis_database.cache.endpoint}:${upstash_redis_database.cache.port}"
}

output "upstash_redis_rest_url" {
  description = "Upstash Redis REST API URL"
  value       = upstash_redis_database.cache.rest_token
  sensitive   = true
}

output "mongodbatlas_connection_string" {
  description = "MongoDB Atlas SRV connection string"
  value       = mongodbatlas_cluster.main.connection_strings[0].standard_srv
  sensitive   = true
}

output "cloudflare_pages_url" {
  description = "Cloudflare Pages deployment URL"
  value       = "https://${cloudflare_pages_project.web.subdomain}"
}

output "render_service_urls" {
  description = "Map of service name to Render public URL"
  value = {
    core-gateway       = module.core_gateway.service_url
    chronos            = module.chronos.service_url
    metis              = module.metis.service_url
    praxis             = module.praxis.service_url
    aegis              = module.aegis.service_url
    orator             = module.orator.service_url
    neuro              = module.neuro.service_url
    argentum           = module.argentum.service_url
    kairos             = module.kairos.service_url
    agent-orchestrator = module.agent_orchestrator.service_url
  }
}
