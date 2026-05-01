resource "render_web_service" "this" {
  name   = var.name
  plan   = var.plan
  region = var.region

  runtime_source = {
    docker = {
      repo_url        = var.github_repo_url
      branch          = var.github_branch
      dockerfile_path = var.dockerfile_path
      context         = var.docker_context
    }
  }

  env_vars = merge(
    { for k, v in var.env_vars : k => { value = v } },
    { for k, v in var.secret_env_vars : k => { value = v } }
  )
}
