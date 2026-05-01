# ---------------------------------------------------------------------------
# Supabase — Postgres
# free-tier: 500 MB storage, 2 projects per org
# ---------------------------------------------------------------------------
resource "supabase_project" "main" {
  organization_id   = var.supabase_org_id
  name              = "ascendos-prod"
  database_password = var.supabase_db_password
  region            = "us-east-1"

  lifecycle {
    ignore_changes = [database_password]
  }
}

# ---------------------------------------------------------------------------
# Upstash — Redis
# free-tier: 10,000 commands/day, 256 MB, 1 database
# ---------------------------------------------------------------------------
resource "upstash_redis_database" "cache" {
  database_name = "ascendos-cache"
  region        = "us-east-1"
  tls           = true
  eviction      = true

  lifecycle {
    prevent_destroy = true
  }
}

# ---------------------------------------------------------------------------
# MongoDB Atlas — M0 free cluster
# free-tier: 512 MB storage, shared CPU/RAM, 1 M0 per project
# ---------------------------------------------------------------------------
resource "mongodbatlas_project" "main" {
  name   = "ascendos"
  org_id = var.mongodbatlas_org_id
}

resource "mongodbatlas_cluster" "main" {
  project_id                  = mongodbatlas_project.main.id
  name                        = "ascendos-free"
  provider_name               = "TENANT"
  backing_provider_name       = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M0" # free-tier: M0 only

  lifecycle {
    prevent_destroy = true
  }
}

resource "mongodbatlas_database_user" "app" {
  project_id         = mongodbatlas_project.main.id
  username           = var.mongodbatlas_db_username
  password           = var.mongodbatlas_db_password
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = "ascendos"
  }
}

# IP allowlist — open for demo; restrict to service IPs before production
resource "mongodbatlas_project_ip_access_list" "any" {
  project_id = mongodbatlas_project.main.id
  cidr_block = "0.0.0.0/0"
  comment    = "demo — restrict before production"
}
