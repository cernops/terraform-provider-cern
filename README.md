# Terraform CERN provider

## Data Sources

### cern_egroup

## Resources

### cern_landb_vm

### cern_landb_vm_card

### cern_landb_vm_interface

## Development and testing

The provider can be built with `go build`. The resulting binary should be place in the following location to match Terraform >= 0.13 requirements: `~/.local/share/terraform/plugins/gitlab.cern.ch/batch-team/cern/1.0.0/linux_amd64/`.