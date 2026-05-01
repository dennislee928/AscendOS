locals {
  github_repo_url = var.github_repo_url
  github_branch   = var.github_branch

  # Supabase connection string — constructed from project outputs
  database_url = "postgresql://postgres:${var.supabase_db_password}@db.${supabase_project.main.id}.supabase.co:5432/postgres"

  # Upstash Redis URL (rediss:// for TLS)
  redis_url = "rediss://:${upstash_redis_database.cache.password}@${upstash_redis_database.cache.endpoint}:${upstash_redis_database.cache.port}"

  # MongoDB Atlas standard connection string
  mongodb_uri = "mongodb+srv://${var.mongodbatlas_db_username}:${var.mongodbatlas_db_password}@${mongodbatlas_cluster.main.connection_strings[0].standard_srv_no_prefix}/ascendos?retryWrites=true&w=majority"

  # Shared data-plane env vars passed to every service
  data_env = {
    DATABASE_URL = local.database_url
    REDIS_URL    = local.redis_url
    MONGODB_URI  = local.mongodb_uri
  }
}

# ---------------------------------------------------------------------------
# Go services
# ---------------------------------------------------------------------------

# free-tier: 750 hr/month per service
module "core_gateway" {
  source = "./modules/render-service"

  name            = "ascendos-core-gateway"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/core-gateway/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT     = "8080"
    NATS_URL = "nats://nats:4222"
  })

  secret_env_vars = {
    CORE_GATEWAY_JWT_SECRET = var.jwt_secret
  }

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "chronos" {
  source = "./modules/render-service"

  name            = "ascendos-chronos"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/chronos/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT = "8080"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "metis" {
  source = "./modules/render-service"

  name            = "ascendos-metis"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/metis/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT = "8081"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "praxis" {
  source = "./modules/render-service"

  name            = "ascendos-praxis"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/praxis/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT = "8082"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# ---------------------------------------------------------------------------
# Rust services
# ---------------------------------------------------------------------------

# free-tier: 750 hr/month per service
module "aegis" {
  source = "./modules/render-service"

  name            = "ascendos-aegis"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/aegis/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT     = "8080"
    RUST_LOG = "info"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "orator" {
  source = "./modules/render-service"

  name            = "ascendos-orator"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/orator/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT     = "8081"
    RUST_LOG = "info"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# ---------------------------------------------------------------------------
# Python / ML services
# ---------------------------------------------------------------------------

# free-tier: 750 hr/month per service
module "neuro" {
  source = "./modules/render-service"

  name            = "ascendos-neuro"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/neuro/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT = "7860"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "argentum" {
  source = "./modules/render-service"

  name            = "ascendos-argentum"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/argentum/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT = "7860"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "kairos" {
  source = "./modules/render-service"

  name            = "ascendos-kairos"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/kairos/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT = "7860"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}

# free-tier: 750 hr/month per service
module "agent_orchestrator" {
  source = "./modules/render-service"

  name            = "ascendos-agent-orchestrator"
  github_repo_url = local.github_repo_url
  github_branch   = local.github_branch
  dockerfile_path = "infra/docker/agent-orchestrator/Dockerfile"
  docker_context  = "."
  health_check_path = "/healthz"

  env_vars = merge(local.data_env, {
    PORT               = "7860"
    ORCHESTRATOR_PORT  = "7860"
  })

  depends_on = [
    supabase_project.main,
    upstash_redis_database.cache,
    mongodbatlas_cluster.main,
  ]
}
