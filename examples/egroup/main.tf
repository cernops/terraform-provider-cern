terraform {
  required_providers {
    cern = {
      source  = "cern-ops/cern"
      version = "> 1.0.0"
    }
  }
}

data "cern_egroup" "batch" {
  name = "batch-3rd"
}

output "batch_members" {
  value = data.cern_egroup.batch.members
}

data "cern_egroup" "ops" {
  name        = "batch-operations"
  query_mails = true
}

output "ops_mails" {
  value = data.cern_egroup.ops.mails
}
