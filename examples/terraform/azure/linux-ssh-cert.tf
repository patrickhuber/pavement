
resource "azurerm_public_ip" "pavement_linux_ssh_cert" {
    name                 = "pavement-public-ip-linux-ssh-cert"
    resource_group_name  = "${azurerm_resource_group.pavement.name}"
    location             = "${azurerm_resource_group.pavement.location}"
    allocation_method    = "Dynamic"
}

resource "azurerm_network_interface" "pavement_linux_ssh_cert" {
    name                = "pavement-nic-linux-ssh-cert"
    resource_group_name = "${azurerm_resource_group.pavement.name}"
    location            = "${azurerm_resource_group.pavement.location}"
    
    network_security_group_id = "${azurerm_network_security_group.pavement.id}"

    ip_configuration{
        name                            = "pavement-private-ip-linux-ssh-cert"
        subnet_id                       = "${azurerm_subnet.pavement.id}"
        private_ip_address_allocation   = "Dynamic"
        public_ip_address_id            = "${azurerm_public_ip.pavement_linux_ssh_cert.id}"
    }
}

resource "azurerm_virtual_machine" "pavement_linux_ssh_cert" {
    name = "pavement-vm-linux-ssh-cert"
    resource_group_name = "${azurerm_resource_group.pavement.name}"
    location = "${azurerm_resource_group.pavement.location}" 
    
    network_interface_ids = ["${azurerm_network_interface.pavement_linux_ssh_cert.id}"]
    vm_size               = "Standard_DS1_v2"

    delete_os_disk_on_termination = true
    
    storage_image_reference {
        publisher = "Canonical"
        offer = "UbuntuServer"
        sku = "18.04-LTS"
        version = "latest"
    }

    storage_os_disk {
        name = "pavement_linux_ssh_os_disk"
        caching = "ReadWrite"
        create_option = "FromImage"
        managed_disk_type = "Standard_LRS"
    }

    os_profile{
        computer_name = "pavement.linux.ssh"
        admin_username = "${var.username}"
    }

    os_profile_linux_config {
        disable_password_authentication = "true"
        ssh_keys {
            key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0QAW7FdVV8fmcwm9Wmnfay9SQbnmElpWrhEbpdPJuboJYE5j6C1PUcAwYOmswAATb4idT8bdl4EvTME4hpiSZZVzblkfojVX9FiNh7Q+tTPjUyBstBOZk6vZ/IVB+wLi2StYFDj4plqAJEAinpAqCMT42dYwRzDBve30jsveH7DPQjRWrQi/mzsihsVkQHxL66vA5zfRJl3gMqJ0BTB1pbEqy4gxLDdxu+OsV+r5s80FU1ASiVPR5REkIdHiFn/ZxbMYIco3UGTWl9xBlD17f+rrWc7voxciqluNBd8BlD6m5QdSY3IkJVGqRmvV6Uu2rZOfk2HRrnr0FoUzDMHxR patrick@DESKTOP-6H368R4"
            path = "/home/${var.username}/.ssh/authorized_keys"
        }
    }
}