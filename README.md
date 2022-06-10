# Terraform CERN provider

This is the Terraform provider to interact with CERN specific services.

For more details about how to use the provider you can check out the `docs/` and
`examples/` directory.

## Development and testing

The provider can be built with `go build`. The resulting binary should be place
in the following location to match Terraform >= 0.13 requirements:
`~/.local/share/terraform/plugins/TODO/$VERSION_HERE/linux_amd64/`.
