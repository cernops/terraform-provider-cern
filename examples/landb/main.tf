resource "cern_landb_vm" "cloud_machine" {
  device_name          = "b7a99n9999"
  location             = {
    building = "0000"
    floor    = "0"
    room     = "0000"
   }
  manufacturer         = "HYPER-V"
  model                = "VIRTUAL MACHINE"
  description          = "Azure Cloud Virtual Machine"
  tag                  = "AZURE CLOUD VM"
  operating_system     = {
    name    = "LINUX"
    version = "UNKNOWN"
  }
  landb_manager_person = {
    name       = "BATCH-XCLOUD-OPERATIONS"
    first_name = "E-GROUP"
    department = "IT"
    group      = "CM"
  }
  responsible_person   = {
    name       = "batch-3rd"
    first_name = "E-GROUP"
    department = "IT"
    group      = "CM"
  }
  user_person          = {
    name       = "batch-3rd"
    first_name = "E-GROUP"
    department = "IT"
    group      = "CM"
  }
  ipv6_ready           = false
  manager_locked       = false
}

resource "cern_landb_vm_card" "cloud_machine_card" {
  vm_name          = cern_landb_vm.cloud_machine.id
  hardware_address = "00-22-48-13-F6-E9"
  card_type        = "Ethernet"
}

resource "cern_landb_vm_interface" "cloud_machine_interface" {
  vm_name              = cern_landb_vm_card.cloud_machine_card.vm_name
  interface_domain     = "cern.ch" # The default
  vm_cluster_name      = "XBATCH-LANDB-AZURE-VM-CLUSTER"
  vm_interface_options = {
    ip           = "188.184.33.10"
    service_name = "S513-C-VM2"
    address_type = "PUBLIC"
  }
}