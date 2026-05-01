variable "github_repo_url" {
  type        = string
  description = "Full HTTPS URL of the GitHub repository"
  default     = "https://github.com/owner/cutting-edge-tools-demo"
}

variable "github_branch" {
  type        = string
  description = "Branch to deploy from"
  default     = "main"
}

# Render
variable "render_api_key" {
  type        = string
  description = "Render API key"
  sensitive   = true
}

variable "render_owner_id" {
  type        = string
  description = "Render owner/user ID (usr-xxxx)"
}

# Cloudflare
variable "cloudflare_api_token" {
  type        = string
  description = "Cloudflare API token with Pages and DNS permissions"
  sensitive   = true
}

variable "cloudflare_account_id" {
  type        = string
  description = "Cloudflare account ID"
}

variable "cloudflare_zone_id" {
  type        = string
  description = "Cloudflare zone ID for custom domain DNS (leave empty if unused)"
  default     = ""
}

# Supabase
variable "supabase_access_token" {
  type        = string
  description = "Supabase personal access token"
  sensitive   = true
}

variable "supabase_org_id" {
  type        = string
  description = "Supabase organization ID"
}

variable "supabase_db_password" {
  type        = string
  description = "Password for the Supabase Postgres database"
  sensitive   = true
}

# Upstash
variable "upstash_email" {
  type        = string
  description = "Upstash account email address"
}

variable "upstash_api_key" {
  type        = string
  description = "Upstash API key"
  sensitive   = true
}

# MongoDB Atlas
variable "mongodbatlas_org_id" {
  type        = string
  description = "MongoDB Atlas organization ID"
}

variable "mongodbatlas_public_key" {
  type        = string
  description = "MongoDB Atlas API public key"
  sensitive   = true
}

variable "mongodbatlas_private_key" {
  type        = string
  description = "MongoDB Atlas API private key"
  sensitive   = true
}

variable "mongodbatlas_db_username" {
  type        = string
  description = "MongoDB Atlas database username for the app"
}

variable "mongodbatlas_db_password" {
  type        = string
  description = "MongoDB Atlas database password for the app"
  sensitive   = true
}

# Application secrets
variable "jwt_secret" {
  type        = string
  description = "JWT signing secret used by core-gateway"
  sensitive   = true
}
