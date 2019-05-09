
resource "azurerm_virtual_network" "pavement" {
    name = "pavement-network"
    resource_group_name = "${azurerm_resource_group.pavement.name}"
    location = "${azurerm_resource_group.pavement.location}"

    address_space = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "pavement" {
    name = "pavement-subnet"
    resource_group_name  = "${azurerm_resource_group.pavement.name}"
    virtual_network_name = "${azurerm_virtual_network.pavement.name}"

    address_prefix = "10.0.0.0/24"
}

resource "azurerm_network_security_group" "pavement" {
    name = "pavement-security-group"    
    resource_group_name = "${azurerm_resource_group.pavement.name}"
    location = "${azurerm_resource_group.pavement.location}"

    security_rule{
        name        = "SSH"
        priority    = 1001
        direction   = "Inbound"
        access      = "Allow"
        protocol    = "Tcp"

        source_address_prefix   = "74.128.138.189"
        source_port_range       = "*"
                
        destination_address_prefix  = "*"
        destination_port_range      = "22"
    }
}