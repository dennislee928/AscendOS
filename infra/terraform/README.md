# AscendOS — Terraform Infrastructure

## Responsibilities

| Tool | What It Manages |
|---|---|
| `deploy-docker-hf.yml` | Builds images, pushes to Docker Hub, deploys to Hugging Face Spaces (primary compute) |
| Terraform (this directory) | Data plane (Supabase, Upstash, MongoDB Atlas), fallback compute (Render), DNS + frontend (Cloudflare) |

Terraform does **not** touch HF Spaces — that is owned entirely by the GitHub Actions workflow.

---

## Free-Tier Breakdown

| Provider | Resource | Free Limit |
|---|---|---|
| Supabase | Postgres | 500 MB storage, 2 projects |
| Upstash | Redis | 10,000 commands/day, 256 MB |
| MongoDB Atlas | M0 cluster | 512 MB storage, shared CPU |
| Render | Web services | 750 hr/month per service, spins down after 15 min idle |
| Cloudflare Pages | Static site | Unlimited requests, 500 builds/month |
| Terraform Cloud | State backend | 500 managed resources |

---

## Prerequisites

- Terraform >= 1.6
- A [Terraform Cloud](https://app.terraform.io) account with org `ascendos` and workspace `ascendos-prod`
- API keys/tokens for every provider (see `terraform.tfvars.example`)

---

## Commands

```bash
# First-time setup
terraform init

# Preview changes
terraform plan

# Apply changes
terraform apply

# Format check
terraform fmt -check -recursive .
```

---

## Backend

State is stored in **Terraform Cloud** (free tier, up to 500 managed resources).

```hcl
terraform {
  cloud {
    organization = "ascendos"
    workspaces { name = "ascendos-prod" }
  }
}
```

Set `TF_API_TOKEN` in your environment (or CI secrets) before running `terraform init`.

---

## Secrets

Never commit real values. Copy `terraform.tfvars.example` to `terraform.tfvars`, fill in secrets, and ensure `terraform.tfvars` is in `.gitignore`.

In CI, secrets are passed as `TF_VAR_*` environment variables — see `.github/workflows/ci-terraform.yml`.
