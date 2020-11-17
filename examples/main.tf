data "cern_egroup" "batch" {
    name = "batch-3rd"
}

output "batch_members" {
    value = data.cern_egroup.batch.members
}
