terraform {
  cloud {
    organization = "ascendos"
    workspaces {
      name = "ascendos-prod"
    }
  }
}
