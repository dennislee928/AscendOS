# free-tier: unlimited requests, 500 builds/month, 1 project
resource "cloudflare_pages_project" "web" {
  account_id        = var.cloudflare_account_id
  name              = "ascendos-web"
  production_branch = "main"

  source = {
    type = "github"
    config = {
      owner                      = split("/", replace(var.github_repo_url, "https://github.com/", ""))[0]
      repo_name                  = split("/", replace(var.github_repo_url, "https://github.com/", ""))[1]
      production_branch          = var.github_branch
      pr_comments_enabled        = true
      deployment_trigger         = "all"
      preview_deployment_setting = "custom"
      preview_branch_includes    = ["dev"]
      preview_branch_excludes    = ["main"]
    }
  }

  build_config = {
    build_command   = "cd apps/web-qwik && npm install && npm run build"
    destination_dir = "apps/web-qwik/dist"
    root_dir        = ""
    build_caching   = true
  }
}
