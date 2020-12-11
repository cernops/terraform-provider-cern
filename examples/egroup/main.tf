terraform {
  required_providers {
    cern = {
      source = "gitlab.cern.ch/batch-team/cern"
      version = "1.0.0"
    }
  }
}

data "cern_egroup" "batch" {
    name = "batch-3rd"
}

output "batch_members" {
    value = data.cern_egroup.batch.members
}
