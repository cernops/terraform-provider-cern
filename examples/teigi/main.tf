terraform {
  required_providers {
    cern = {
      source = "gitlab.cern.ch/batch-team/cern"
      version = "1.0.0"
    }
  }
}

data "cern_teigi_secret" "oops" {
  key       = "test_this"
  hostgroup = "playground"
}

output "secret" {
    value = data.cern_teigi_secret.oops.secret
}
