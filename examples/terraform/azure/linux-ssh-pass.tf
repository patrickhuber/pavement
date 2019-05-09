
resource "azurerm_public_ip" "pavement_linux_ssh_pass" {
    name                 = "pavement-public-ip-linux-ssh-pass"
    resource_group_name  = "${azurerm_resource_group.pavement.name}"
    location             = "${azurerm_resource_group.pavement.location}"
    allocation_method    = "Dynamic"
}

resource "azurerm_network_interface" "pavement_linux_ssh_pass" {
    name                = "pavement-nic-linux-ssh-pass"
    resource_group_name = "${azurerm_resource_group.pavement.name}"
    location            = "${azurerm_resource_group.pavement.location}"
    
    network_security_group_id = "${azurerm_network_security_group.pavement.id}"

    ip_configuration{
        name                            = "pavement-private-ip-linux-ssh-pass"
        subnet_id                       = "${azurerm_subnet.pavement.id}"
        private_ip_address_allocation   = "Dynamic"
        public_ip_address_id            = "${azurerm_public_ip.pavement_linux_ssh_pass.id}"
    }
}

resource "azurerm_virtual_machine" "pavement_linux_ssh_pass" {
    name = "pavement-vm-linux-ssh-pass"
    resource_group_name = "${azurerm_resource_group.pavement.name}"
    location = "${azurerm_resource_group.pavement.location}" 
    
    network_interface_ids = ["${azurerm_network_interface.pavement_linux_ssh_pass.id}"]
    vm_size               = "Standard_DS1_v2"

    delete_os_disk_on_termination = true
    
    storage_image_reference {
        publisher = "Canonical"
        offer = "UbuntuServer"
        sku = "18.04-LTS"
        version = "latest"
    }

    storage_os_disk {
        name = "pavement_linux_ssh_pass_os_disk"
        caching = "ReadWrite"
        create_option = "FromImage"
        managed_disk_type = "Standard_LRS"
    }

    os_profile{
        computer_name = "pavementlinuxsshpass"
        admin_username = "${var.username}"
        admin_password = "${var.password}"
    }

    os_profile_linux_config {
        disable_password_authentication = "false"
    }
}