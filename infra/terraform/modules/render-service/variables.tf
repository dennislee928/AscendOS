variable "name" {
  type        = string
  description = "Name of the Render web service"
}

variable "plan" {
  type        = string
  description = "Render service plan"
  default     = "free" # free-tier: 750 hr/month, spins down after 15 min idle
}

variable "region" {
  type        = string
  description = "Render deployment region"
  default     = "oregon"
}

variable "github_repo_url" {
  type        = string
  description = "Full HTTPS URL of the GitHub repository"
}

variable "github_branch" {
  type        = string
  description = "Branch to deploy from"
  default     = "main"
}

variable "dockerfile_path" {
  type        = string
  description = "Path to the Dockerfile relative to the repo root"
}

variable "docker_context" {
  type        = string
  description = "Docker build context path relative to the repo root"
  default     = "."
}

variable "health_check_path" {
  type        = string
  description = "HTTP path used by Render for health checks"
  default     = "/healthz"
}

variable "env_vars" {
  type        = map(string)
  description = "Non-secret environment variables"
  default     = {}
}

variable "secret_env_vars" {
  type        = map(string)
  description = "Secret environment variables (stored encrypted by Render)"
  default     = {}
  sensitive   = true
}
