terraform {
  required_version = ">= 1.6"

  required_providers {
    render = {
      source  = "render-oss/render"
      version = "~> 1.5"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.36"
    }
    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.4"
    }
    upstash = {
      source  = "upstash/upstash"
      version = "~> 1.5"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.16"
    }
  }
}

provider "render" {
  api_key  = var.render_api_key
  owner_id = var.render_owner_id
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

provider "supabase" {
  access_token = var.supabase_access_token
}

provider "upstash" {
  email   = var.upstash_email
  api_key = var.upstash_api_key
}

provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
