# Terraform CERN provider

## Provider configuration

The provider configuration block accepts the following arguments:

* `ldap_server` (_optional_, `ldap://xldap.cern.ch:389`): ldap URL to use for queries by the e-group data source. It can be set with the environment variable `CERN_LDAP_SERVER`.

* `landb_endpoint` (_optional_, default: `https://network.cern.ch/sc/soap/soap.fcgi?v=6`): LanDB URL to interact with to create resources. It can be set with the environment variable `CERN_LANDB_ENDPOINT`.

* `landb_username` (_optional_): Username with privileges to interact with LanDB. It is recommended to set this field using the environment variable `CERN_LANDB_USERNAME`.

* `landb_password` (_optional_): Password to interact with LanDB. It is recommended to set this field using the environment variable `CERN_LANDB_PASSWORD`.

* `teigi_endpoint` (_optional_, default: `https://woger.cern.ch:8201`): It can be set with the environment variable `CERN_TEIGI_ENDPOINT`.

Example usage:

```
provider "cern" {
   landb_username = "svc.account"
   teigi_endpoint = "https://teigi-xyz.cern.ch:8201"
}
```

## Data Sources

### cern_egroup

### cern_teigi_secret

Reads a secret from Teigi with specific hostgroup and key.

#### Example usage

```hcl
data "cern_teigi_secret" "data" {
  key       = "gitlab_service_token"
  hostgroup = "awesome/hostgroup"
}

resource "null_resource" "dummy" {
  provisioner "local-exec" {
    command = "echo \"Our token is: ${data.cern_teigi_secret.data.secret}\""
  }
}
```
#### Argument reference

The following arguments are supported:

* __key__ - (Required) The key value to search the secret for.
* __hostgroup__ - (Required) The hostgroup to query for the secret.

#### Attributes reference

The following attributes are exported:

* __secret__ - A string containing the secret retrieved from Teigi.

## Resources

### cern_landb_vm

### cern_landb_vm_card

### cern_landb_vm_interface

## Development and testing

The provider can be built with `go build`. The resulting binary should be place in the following location to match Terraform >= 0.13 requirements: `~/.local/share/terraform/plugins/gitlab.cern.ch/batch-team/cern/1.0.0/linux_amd64/`.
